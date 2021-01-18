package campaign

import (
	"strings"
)

type CampaignFormatter struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	ImageURL         string `json:"image_url"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
	Slug             string `json:"slug"`
}

type CampaignDetailFormatter struct {
	ID               int                     `json:"id"`
	Name             string                  `json:"name"`
	ShortDescription string                  `json:"short_description"`
	Description      string                  `json:"description"`
	ImageURL         string                  `json:"image_url"`
	GoalAmount       int                     `json:"goal_amount"`
	CurrentAmount    int                     `json:"current_amount"`
	BackerCount      int                     `json:"backer_count"`
	UserID           int                     `json:"user_id"`
	Slug             string                  `json:"slug"`
	Perks            []string                `json:"perks"`
	User             CampignUserFormatter    `json:"user"`
	Images           []CampignImageFormatter `json:"images"`
}

type CampignUserFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}
type CampignImageFormatter struct {
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

func FormatCampaign(campaign Campaign) CampaignFormatter {

	campaignFormatter := CampaignFormatter{}
	campaignFormatter.ID = campaign.ID
	campaignFormatter.UserID = campaign.UserID
	campaignFormatter.Name = campaign.Name
	campaignFormatter.ShortDescription = campaign.ShortDescription
	campaignFormatter.GoalAmount = campaign.GoalAmount
	campaignFormatter.CurrentAmount = campaign.CurrentAmount

	campaignFormatter.Slug = campaign.Slug
	campaignFormatter.ImageURL = ""

	if len(campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	return campaignFormatter
}

func FormatCampaigns(campaigns []Campaign) []CampaignFormatter {

	allFormatter := []CampaignFormatter{}
	for _, campaign := range campaigns {

		campaignFormatter := FormatCampaign(campaign)

		allFormatter = append(allFormatter, campaignFormatter)
	}

	return allFormatter
}

func FormatCampaignDetail(campaign Campaign) CampaignDetailFormatter {

	data := CampaignDetailFormatter{}

	data.ID = campaign.ID
	data.Name = campaign.Name
	data.ShortDescription = campaign.ShortDescription
	data.Description = campaign.Description
	data.GoalAmount = campaign.GoalAmount
	data.CurrentAmount = campaign.CurrentAmount
	data.BackerCount = campaign.BackerCount
	data.Slug = campaign.Slug
	data.ImageURL = ""

	if len(campaign.CampaignImages) > 0 {
		data.ImageURL = campaign.CampaignImages[0].FileName
	}
	var perks []string
	for _, perk := range strings.Split(campaign.Perks, ",") {

		perks = append(perks, strings.TrimSpace(perk))
	}

	data.Perks = perks

	user := campaign.User

	campaignUser := CampignUserFormatter{}
	campaignUser.Name = user.Name
	campaignUser.ImageURL = user.AvatarFileName

	data.User = campaignUser

	images := []CampignImageFormatter{}

	for _, img := range campaign.CampaignImages {

		campaignImage := CampignImageFormatter{}

		campaignImage.ImageURL = img.FileName

		isPrimary := false
		if img.IsPrimary == 1 {
			isPrimary = true
		}

		campaignImage.IsPrimary = isPrimary

		images = append(images, campaignImage)

	}

	data.Images = images
	return data
}
