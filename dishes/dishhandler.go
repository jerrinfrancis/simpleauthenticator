package dishes

import (
	"encoding/json"
	"net/http"
)

func GetDishes(w http.ResponseWriter, r *http.Request) {

	dishes := []Dish{{
		Name:          "Mutton Biriyani",
		ChefName:      "Lincy T Elizabeth",
		Price:         "200",
		ThumbImageURL: "https://www.cubesnjuliennes.com/wp-content/uploads/2020/06/Mutton-Biryani-Recipe.jpg",
	},
		{
			Name:          "Chicken Biriyani",
			ChefName:      "Jerrin Francis",
			Price:         "150",
			ThumbImageURL: "https://www.indianhealthyrecipes.com/wp-content/uploads/2019/02/chicken-biryani-recipe.jpg",
		},
		{
			Name:          "Vanilla Icecream",
			ChefName:      "Princy T Elizabeth",
			Price:         "150",
			ThumbImageURL: "https://www.indianhealthyrecipes.com/wp-content/uploads/2016/05/vanilla-ice-cream-recipe-500x375.jpg",
		},
		{
			Name:          "Valsan",
			ChefName:      "Thankachan O",
			Price:         "10",
			ThumbImageURL: "https://i.ytimg.com/vi/XmjcVuhYXB8/hqdefault.jpg",
		}}
	a, _ := json.Marshal(&dishes)
	w.WriteHeader(http.StatusOK)
	w.Write(a)

}
