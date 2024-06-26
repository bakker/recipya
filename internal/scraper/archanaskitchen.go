package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeArchanasKitchen(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	description := root.Find("span[itemprop='description']").Text()
	description = strings.TrimPrefix(description, "\n")
	description = strings.ReplaceAll(description, "\u00a0", " ")
	rs.Description.Value = strings.TrimSpace(description)

	image, _ := root.Find("img[itemprop='image']").Attr("src")
	rs.Image.Value = "https://www.archanaskitchen.com" + image

	root.Find("li[itemprop='keywords'] a").Each(func(_ int, s *goquery.Selection) {
		rs.Keywords.Values += strings.TrimSpace(s.Text()) + ","
	})
	rs.Keywords.Values = strings.TrimSuffix(rs.Keywords.Values, ",")

	nodes := root.Find("li[itemprop='ingredients']")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		v := strings.ReplaceAll(s.Text(), "\n", "")
		v = strings.ReplaceAll(v, "\t", "")
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(strings.ReplaceAll(v, " , ", ", ")))
	})

	nodes = root.Find("li[itemprop='recipeInstructions'] p")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		v := strings.ReplaceAll(s.Text(), "\u00a0", " ")
		v = strings.TrimSpace(strings.ReplaceAll(v, " .", "."))
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(v))
	})

	rs.PrepTime, _ = root.Find("span[itemprop='prepTime']").Attr("content")
	rs.CookTime, _ = root.Find("span[itemprop='cookTime']").Attr("content")
	rs.DatePublished, _ = root.Find("span[itemprop='datePublished']").Attr("content")
	rs.DateModified, _ = root.Find("span[itemprop='dateModified']").Attr("content")
	rs.Yield.Value = findYield(root.Find("span[itemprop='recipeYield'] p").Text())
	rs.Name = root.Find("h1[itemprop='name']").Text()
	rs.Category = &models.Category{Value: root.Find(".recipeCategory a").Text()}
	rs.Cuisine = &models.Cuisine{Value: root.Find("span[itemprop='recipeCuisine'] a").Text()}

	return rs, nil
}
