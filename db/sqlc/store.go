package db

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Store interface {
	Querier
	Delivery(ctx context.Context, arg DeliveryParams) ([]DeliveryResult, error)
	CreateCampaign(tx context.Context, arg CreateCampaignParams) (CreateCampaignResult, error)
	ReadCampaign(ctx context.Context, cid string) (CompleteCampaign, error)
	ToggleStatus(ctx context.Context, cid string) error
	UpdateCampaignName(ctx context.Context, arg UpdateCampaignNameParams) (Campaign, error)
	UpdateCampaignCta(ctx context.Context, arg UpdateCampaignCtaParams) (Campaign, error)
	UpdateCampaignImage(ctx context.Context, arg UpdateCampaignImageParams) (Campaign, error)
	UpdateTargetApp(ctx context.Context, arg UpdateTargetAppParams) (TargetApp, error)
	UpdateTargetCountry(ctx context.Context, arg UpdateTargetCountryParams) (TargetCountry, error)
	UpdateTargetOs(ctx context.Context, arg UpdateTargetOsParams) (TargetOs, error)
}

type SQLStore struct {
	*Queries
	db      *pgxpool.Pool
	rClient *redis.Client
}

func NewStore(db *pgxpool.Pool, rClient *redis.Client) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
		rClient: rClient,
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

func csvToSlice(csv string) []string {
	list := strings.Split(csv, ",")
	for i, item := range list {
		list[i] = strings.TrimSpace(item)
	}
	return list
}

func contains(slice []string, str string) bool {
	str = strings.ToLower(str)
	for _, v := range slice {
		if strings.ToLower(v) == str {
			return true
		}
	}
	return false
}

func shouldInclude(csv string, arg string, rule string) bool {
	list := csvToSlice(csv)
	contains := contains(list, arg)

	if (rule == "include" && !contains) || (rule == "exclude" && contains) {
		return false
	}
	return true
}

type DeliveryParams struct {
	AppID   string `json:"app_id"`
	Country string `json:"country"`
	Os      string `json:"os"`
}

type DeliveryResult struct {
	Cid string `json:"cid"`
	Img string `json:"img"`
	Cta string `json:"cta"`
}

// Helper functions for caching query results

func (store *SQLStore) getCachedActiveCampaigns(ctx context.Context) ([]Campaign, error) {
	cacheKey := "active_campaigns"
	var campaigns []Campaign

	cachedData, err := store.rClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		err = json.Unmarshal(cachedData, &campaigns)
		if err == nil {
			return campaigns, nil
		}
		fmt.Printf("Redis JSON unmarshal error for active campaigns: %v\n", err)
	} else if err != redis.Nil {
		fmt.Printf("Redis Get error for active campaigns: %v\n", err)
	}

	campaigns, err = store.ListActiveCampaigns(ctx)
	if err != nil {
		return nil, err
	}

	if len(campaigns) > 0 {
		jsonData, err := json.Marshal(campaigns)
		if err != nil {
			fmt.Printf("JSON marshal error for active campaigns: %v\n", err)
		} else {
			err = store.rClient.Set(ctx, cacheKey, jsonData, 15*time.Minute).Err()
			if err != nil {
				fmt.Printf("Redis Set error for active campaigns: %v\n", err)
			}
		}
	}

	return campaigns, nil
}

func (store *SQLStore) getCachedTargetApp(ctx context.Context, cid string) (*TargetApp, error) {
	cacheKey := fmt.Sprintf("target_app:%s", cid)
	var targetApp TargetApp

	cachedData, err := store.rClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		err = json.Unmarshal(cachedData, &targetApp)
		if err == nil {
			return &targetApp, nil
		}
		fmt.Printf("Redis JSON unmarshal error for target app: %v\n", err)
	} else if err != redis.Nil {
		fmt.Printf("Redis Get error for target app: %v\n", err)
	}

	result, err := store.GetTargetApp(ctx, cid)
	if err != nil {
		if err == pgx.ErrNoRows {
			store.rClient.Set(ctx, cacheKey, "null", 15*time.Minute)
			return nil, err
		}
		return nil, err
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		fmt.Printf("JSON marshal error for target app: %v\n", err)
	} else {
		err = store.rClient.Set(ctx, cacheKey, jsonData, 15*time.Minute).Err()
		if err != nil {
			fmt.Printf("Redis Set error for target app: %v\n", err)
		}
	}

	return &result, nil
}

func (store *SQLStore) getCachedTargetCountry(ctx context.Context, cid string) (*TargetCountry, error) {
	cacheKey := fmt.Sprintf("target_country:%s", cid)
	var targetCountry TargetCountry

	cachedData, err := store.rClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		err = json.Unmarshal(cachedData, &targetCountry)
		if err == nil {
			return &targetCountry, nil
		}
		fmt.Printf("Redis JSON unmarshal error for target country: %v\n", err)
	} else if err != redis.Nil {
		fmt.Printf("Redis Get error for target country: %v\n", err)
	}

	result, err := store.GetTargetCountry(ctx, cid)
	if err != nil {
		if err == pgx.ErrNoRows {
			store.rClient.Set(ctx, cacheKey, "null", 15*time.Minute)
			return nil, err
		}
		return nil, err
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		fmt.Printf("JSON marshal error for target country: %v\n", err)
	} else {
		err = store.rClient.Set(ctx, cacheKey, jsonData, 15*time.Minute).Err()
		if err != nil {
			fmt.Printf("Redis Set error for target country: %v\n", err)
		}
	}

	return &result, nil
}

func (store *SQLStore) getCachedTargetOs(ctx context.Context, cid string) (*TargetOs, error) {
	cacheKey := fmt.Sprintf("target_os:%s", cid)
	var targetOs TargetOs

	cachedData, err := store.rClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		err = json.Unmarshal(cachedData, &targetOs)
		if err == nil {
			return &targetOs, nil
		}
		fmt.Printf("Redis JSON unmarshal error for target os: %v\n", err)
	} else if err != redis.Nil {
		fmt.Printf("Redis Get error for target os: %v\n", err)
	}

	result, err := store.GetTargetOs(ctx, cid)
	if err != nil {
		if err == pgx.ErrNoRows {
			store.rClient.Set(ctx, cacheKey, "null", 15*time.Minute)
			return nil, err
		}
		return nil, err
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		fmt.Printf("JSON marshal error for target os: %v\n", err)
	} else {
		err = store.rClient.Set(ctx, cacheKey, jsonData, 15*time.Minute).Err()
		if err != nil {
			fmt.Printf("Redis Set error for target os: %v\n", err)
		}
	}

	return &result, nil
}

func (store *SQLStore) Delivery(ctx context.Context, arg DeliveryParams) ([]DeliveryResult, error) {
	cacheKey := fmt.Sprintf("delivery:%s:%s:%s", arg.AppID, arg.Country, arg.Os)
	var results []DeliveryResult

	cachedData, err := store.rClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		err = json.Unmarshal(cachedData, &results)
		if err == nil {
			return results, nil
		}
		fmt.Printf("Redis JSON unmarshal error: %v\n", err)
	} else if err != redis.Nil {
		fmt.Printf("Redis Get error: %v\n", err)
	}

	active_campaigns, err := store.getCachedActiveCampaigns(ctx)
	if err != nil {
		return []DeliveryResult{}, err
	}

	var campaigns []Campaign
	for _, campaign := range active_campaigns {
		target_app, err := store.getCachedTargetApp(ctx, campaign.Cid)
		if err != nil && err != pgx.ErrNoRows {
			return []DeliveryResult{}, err
		}
		if err != pgx.ErrNoRows {
			if !shouldInclude(target_app.AppID, arg.AppID, string(target_app.Rule)) {
				continue
			}
		}

		target_country, err := store.getCachedTargetCountry(ctx, campaign.Cid)
		if err != nil && err != pgx.ErrNoRows {
			return []DeliveryResult{}, err
		}
		if err != pgx.ErrNoRows {
			if !shouldInclude(target_country.Country, arg.Country, string(target_country.Rule)) {
				continue
			}
		}

		target_os, err := store.getCachedTargetOs(ctx, campaign.Cid)
		if err != nil && err != pgx.ErrNoRows {
			return []DeliveryResult{}, err
		}
		if err != pgx.ErrNoRows {
			if !shouldInclude(target_os.Os, arg.Os, string(target_os.Rule)) {
				continue
			}
		}

		campaigns = append(campaigns, campaign)
	}

	var result []DeliveryResult
	for _, campaign := range campaigns {
		result = append(result, DeliveryResult{
			Cid: campaign.Cid,
			Img: campaign.Img,
			Cta: campaign.Cta,
		})
	}

	if len(result) > 0 {
		jsonData, err := json.Marshal(result)
		if err != nil {
			fmt.Printf("JSON marshal error: %v\n", err)
		} else {
			err = store.rClient.Set(ctx, cacheKey, jsonData, 30*time.Minute).Err()
			if err != nil {
				fmt.Printf("Redis Set error: %v\n", err)
			}
		}
	}

	return result, nil
}

type CreateCampaignParams struct {
	Cid         string   `json:"cid"`
	Name        string   `json:"name"`
	Img         string   `json:"img"`
	Cta         string   `json:"cta"`
	AppID       string   `json:"app_id"`
	AppRule     RuleType `json:"app_rule"`
	Country     string   `json:"country"`
	CountryRule RuleType `json:"country_rule"`
	Os          string   `json:"os"`
	OsRule      RuleType `json:"os_rule"`
}

type CreateCampaignResult struct {
	Cid         string     `json:"cid"`
	Name        string     `json:"name"`
	Img         string     `json:"img"`
	Cta         string     `json:"cta"`
	AppID       string     `json:"app_id"`
	AppRule     RuleType   `json:"app_rule"`
	Country     string     `json:"country"`
	CountryRule RuleType   `json:"country_rule"`
	Os          string     `json:"os"`
	OsRule      RuleType   `json:"os_rule"`
	Status      StatusType `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (store *SQLStore) CreateCampaign(ctx context.Context, arg CreateCampaignParams) (CreateCampaignResult, error) {
	var result CreateCampaignResult

	err := store.execTx(ctx, func(q *Queries) error {
		campaign, err := q.AddCampaign(ctx, AddCampaignParams{
			Cid:  arg.Cid,
			Name: arg.Name,
			Img:  arg.Img,
			Cta:  arg.Cta,
		})
		if err != nil {
			return err
		}

		result = CreateCampaignResult{
			Cid:       campaign.Cid,
			Name:      campaign.Name,
			Img:       campaign.Img,
			Cta:       campaign.Cta,
			Status:    campaign.Status,
			CreatedAt: campaign.CreatedAt,
		}

		if arg.AppID != "" {
			targetApp, err := q.AddTargetApp(ctx, AddTargetAppParams{
				Cid:   arg.Cid,
				AppID: arg.AppID,
				Rule:  arg.AppRule,
			})
			if err != nil {
				return err
			}
			result.AppID = targetApp.AppID
			result.AppRule = targetApp.Rule
		}

		if arg.Country != "" {
			targetCountry, err := q.AddTargetCountry(ctx, AddTargetCountryParams{
				Cid:     arg.Cid,
				Country: arg.Country,
				Rule:    arg.CountryRule,
			})
			if err != nil {
				return err
			}
			result.Country = targetCountry.Country
			result.CountryRule = targetCountry.Rule
		}

		if arg.Os != "" {
			targetOs, err := q.AddTargetOs(ctx, AddTargetOsParams{
				Cid:  arg.Cid,
				Os:   arg.Os,
				Rule: arg.OsRule,
			})
			if err != nil {
				return err
			}
			result.Os = targetOs.Os
			result.OsRule = targetOs.Rule
		}

		return nil
	})

	return result, err
}

type CompleteCampaign struct {
	Cid         string     `json:"cid"`
	Name        string     `json:"name"`
	Img         string     `json:"img"`
	Cta         string     `json:"cta"`
	AppID       string     `json:"app_id"`
	AppRule     RuleType   `json:"app_rule"`
	Country     string     `json:"country"`
	CountryRule RuleType   `json:"country_rule"`
	Os          string     `json:"os"`
	OsRule      RuleType   `json:"os_rule"`
	Status      StatusType `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (store *SQLStore) ReadCampaign(ctx context.Context, cid string) (CompleteCampaign, error) {
	campaign, err := store.GetCampaign(ctx, cid)
	if err != nil {
		return CompleteCampaign{}, err
	}

	TargetApp, _ := store.GetTargetApp(ctx, cid)
	TargetCountry, _ := store.GetTargetCountry(ctx, cid)
	TargetOs, _ := store.GetTargetOs(ctx, cid)

	return CompleteCampaign{
		Cid:         cid,
		Name:        campaign.Name,
		Img:         campaign.Img,
		Cta:         campaign.Cta,
		AppID:       TargetApp.AppID,
		AppRule:     TargetApp.Rule,
		Country:     TargetCountry.Country,
		CountryRule: TargetCountry.Rule,
		Os:          TargetOs.Os,
		OsRule:      TargetOs.Rule,
		Status:      campaign.Status,
		CreatedAt:   campaign.CreatedAt,
	}, nil
}

func (store *SQLStore) createHistory(ctx context.Context, q *Queries, args []createCampaignHistoryParams) error {
	for _, arg := range args {
		if arg.OldValue != arg.NewValue {
			err := q.createCampaignHistory(ctx, arg)
			if err != nil {
				return fmt.Errorf("failed to create history for %s: %v", arg.FieldChanged, err)
			}
		}
	}
	return nil
}

func (store *SQLStore) ToggleStatus(ctx context.Context, cid string) error {
	return store.execTx(ctx, func(q *Queries) error {
		campaign, err := q.GetCampaign(ctx, cid)
		if err != nil {
			return err
		}

		status, err := q.toggleStatus(ctx, cid)
		if err != nil {
			return err
		}

		return store.createHistory(ctx, q, []createCampaignHistoryParams{
			{
				Cid:          cid,
				FieldChanged: "status",
				OldValue:     string(campaign.Status),
				NewValue:     string(status),
			},
		})
	})
}

type UpdateCampaignNameParams struct {
	Cid  string `json:"cid"`
	Name string `json:"name"`
}

func (store *SQLStore) UpdateCampaignName(ctx context.Context, arg UpdateCampaignNameParams) (Campaign, error) {
	var campaign Campaign
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		oldCampaign, err := q.GetCampaign(ctx, arg.Cid)
		if err != nil {
			return err
		}

		campaign, err = q.updateCampaignName(ctx, updateCampaignNameParams{
			Cid:  arg.Cid,
			Name: arg.Name,
		})
		if err != nil {
			return err
		}

		return store.createHistory(ctx, q, []createCampaignHistoryParams{
			{
				Cid:          arg.Cid,
				FieldChanged: "name",
				OldValue:     oldCampaign.Name,
				NewValue:     campaign.Name,
			},
		})
	})
	return campaign, err
}

type UpdateCampaignCtaParams struct {
	Cid string `json:"cid"`
	Cta string `json:"cta"`
}

func (store *SQLStore) UpdateCampaignCta(ctx context.Context, arg UpdateCampaignCtaParams) (Campaign, error) {
	var campaign Campaign
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		oldCampaign, err := q.GetCampaign(ctx, arg.Cid)
		if err != nil {
			return err
		}

		campaign, err = q.updateCampaignCta(ctx, updateCampaignCtaParams{
			Cid: arg.Cid,
			Cta: arg.Cta,
		})
		if err != nil {
			return err
		}

		return store.createHistory(ctx, q, []createCampaignHistoryParams{
			{
				Cid:          arg.Cid,
				FieldChanged: "cta",
				OldValue:     oldCampaign.Cta,
				NewValue:     campaign.Cta,
			},
		})
	})
	return campaign, err
}

type UpdateCampaignImageParams struct {
	Cid string `json:"cid"`
	Img string `json:"img"`
}

func (store *SQLStore) UpdateCampaignImage(ctx context.Context, arg UpdateCampaignImageParams) (Campaign, error) {
	var campaign Campaign
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		oldCampaign, err := q.GetCampaign(ctx, arg.Cid)
		if err != nil {
			return err
		}

		campaign, err = q.updateCampaignImage(ctx, updateCampaignImageParams{
			Cid: arg.Cid,
			Img: arg.Img,
		})
		if err != nil {
			return err
		}

		return store.createHistory(ctx, q, []createCampaignHistoryParams{
			{
				Cid:          arg.Cid,
				FieldChanged: "img",
				OldValue:     oldCampaign.Img,
				NewValue:     campaign.Img,
			},
		})
	})
	return campaign, err
}

type UpdateTargetAppParams struct {
	Cid   string   `json:"cid"`
	AppID string   `json:"app_id"`
	Rule  RuleType `json:"rule"`
}

func (store *SQLStore) UpdateTargetApp(ctx context.Context, arg UpdateTargetAppParams) (TargetApp, error) {
	var targetApp TargetApp
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		oldTarget, err := q.GetTargetApp(ctx, arg.Cid)
		if err != nil {
			return err
		}

		targetApp, err = q.updateTargetApp(ctx, updateTargetAppParams{
			Cid:   arg.Cid,
			AppID: arg.AppID,
			Rule:  arg.Rule,
		})
		if err != nil {
			return err
		}

		return store.createHistory(ctx, q, []createCampaignHistoryParams{
			{
				Cid:          arg.Cid,
				FieldChanged: "app_id",
				OldValue:     oldTarget.AppID,
				NewValue:     targetApp.AppID,
			},
			{
				Cid:          arg.Cid,
				FieldChanged: "app_rule",
				OldValue:     string(oldTarget.Rule),
				NewValue:     string(targetApp.Rule),
			},
		})
	})
	return targetApp, err
}

type UpdateTargetOsParams struct {
	Cid  string   `json:"cid"`
	Os   string   `json:"os"`
	Rule RuleType `json:"rule"`
}

func (store *SQLStore) UpdateTargetOs(ctx context.Context, arg UpdateTargetOsParams) (TargetOs, error) {
	var targetOs TargetOs
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		oldTarget, err := q.GetTargetOs(ctx, arg.Cid)
		if err != nil {
			return err
		}

		targetOs, err = q.updateTargetOs(ctx, updateTargetOsParams{
			Cid:  arg.Cid,
			Os:   arg.Os,
			Rule: arg.Rule,
		})
		if err != nil {
			return err
		}

		return store.createHistory(ctx, q, []createCampaignHistoryParams{
			{
				Cid:          arg.Cid,
				FieldChanged: "os",
				OldValue:     oldTarget.Os,
				NewValue:     targetOs.Os,
			},
			{
				Cid:          arg.Cid,
				FieldChanged: "os_rule",
				OldValue:     string(oldTarget.Rule),
				NewValue:     string(targetOs.Rule),
			},
		})
	})
	return targetOs, err
}

type UpdateTargetCountryParams struct {
	Cid     string   `json:"cid"`
	Country string   `json:"country"`
	Rule    RuleType `json:"rule"`
}

func (store *SQLStore) UpdateTargetCountry(ctx context.Context, arg UpdateTargetCountryParams) (TargetCountry, error) {
	var targetCountry TargetCountry
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		oldTarget, err := q.GetTargetCountry(ctx, arg.Cid)
		if err != nil {
			return err
		}

		targetCountry, err = q.updateTargetCountry(ctx, updateTargetCountryParams{
			Cid:     arg.Cid,
			Country: arg.Country,
			Rule:    arg.Rule,
		})
		if err != nil {
			return err
		}

		return store.createHistory(ctx, q, []createCampaignHistoryParams{
			{
				Cid:          arg.Cid,
				FieldChanged: "country",
				OldValue:     oldTarget.Country,
				NewValue:     targetCountry.Country,
			},
			{
				Cid:          arg.Cid,
				FieldChanged: "country_rule",
				OldValue:     string(oldTarget.Rule),
				NewValue:     string(targetCountry.Rule),
			},
		})
	})
	return targetCountry, err
}
