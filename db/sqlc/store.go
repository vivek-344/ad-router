package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	Delivery(ctx context.Context, arg DeliveryParams) ([]DeliveryResult, error)
	AddCampaign(tx context.Context, arg AddCampaignParams) (AddCampaignResult, error)
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
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
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

func (store *SQLStore) Delivery(ctx context.Context, arg DeliveryParams) ([]DeliveryResult, error) {
	var campaigns []Campaign
	var result []DeliveryResult
	active_campaigns, err := store.ListAllActiveCampaign(ctx)
	if err != nil {
		return []DeliveryResult{}, err
	}

	for _, campaign := range active_campaigns {
		target_app, err := store.GetTargetApp(ctx, campaign.Cid)
		if err != nil && err != pgx.ErrNoRows {
			return []DeliveryResult{}, err
		}
		if err != pgx.ErrNoRows {
			if !shouldInclude(target_app.AppID, arg.AppID, string(target_app.Rule)) {
				continue
			}
		}

		target_country, err := store.GetTargetCountry(ctx, campaign.Cid)
		if err != nil && err != pgx.ErrNoRows {
			return []DeliveryResult{}, err
		}
		if err != pgx.ErrNoRows {
			if !shouldInclude(target_country.Country, arg.Country, string(target_country.Rule)) {
				continue
			}
		}

		target_os, err := store.GetTargetOs(ctx, campaign.Cid)
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

	for _, campaign := range campaigns {
		result = append(result, DeliveryResult{
			Cid: campaign.Cid,
			Img: campaign.Img,
			Cta: campaign.Cta,
		})
	}

	return result, nil
}

type AddCampaignParams struct {
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

type AddCampaignResult struct {
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

func (store *SQLStore) AddCampaign(ctx context.Context, arg AddCampaignParams) (AddCampaignResult, error) {
	var result AddCampaignResult

	err := store.execTx(ctx, func(q *Queries) error {
		campaign, err := q.CreateCampaign(ctx, CreateCampaignParams{
			Cid:  arg.Cid,
			Name: arg.Name,
			Img:  arg.Img,
			Cta:  arg.Cta,
		})
		if err != nil {
			return err
		}

		result = AddCampaignResult{
			Cid:       campaign.Cid,
			Name:      campaign.Name,
			Img:       campaign.Img,
			Cta:       campaign.Cta,
			Status:    campaign.Status,
			CreatedAt: campaign.CreatedAt,
		}

		if arg.AppID != "" {
			targetApp, err := q.CreateTargetApp(ctx, CreateTargetAppParams{
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
			targetCountry, err := q.CreateTargetCountry(ctx, CreateTargetCountryParams{
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
			targetOs, err := q.CreateTargetOs(ctx, CreateTargetOsParams{
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
