package cmd

import (
	"dj/pkg/dj"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

var (
	image     bool
	slack     bool
	imagePath string
)

var rootCmd = &cobra.Command{
	Use:   "dj",
	Short: "dj retrieves dad jokes",
	Long:  "dj is a command line tool for retrieving dad jokes from icanhazdadjoke.com",
	Run: func(cmd *cobra.Command, args []string) {
		baseUrl, err := url.Parse("https://icanhazdadjoke.com")
		if err != nil {
			panic(err)
		}

		c := dj.Client{
			BaseUrl:    baseUrl,
			UserAgent:  "dj-go (https://github.com/rchaganti/dadjoke-go)",
			HttpClient: &http.Client{},
		}

		if slack {
			j, err := c.GetJokeAsSlackMessage()
			if err != nil {
				panic(err)
			}
			json, err := json.MarshalIndent(j, "", "  ")
			if err != nil {
				panic(err)
			}

			fmt.Println(string(json))
			return
		} else {
			j, err := c.GetJoke()
			if err != nil {
				panic(err)
			}

			if image {
				jid := j.ID
				err := c.GetJokeAsImage(jid, imagePath)
				if err != nil {
					panic(err)
				}
			} else {
				fmt.Println(j.Joke)
			}
		}
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&image, "asImage", "i", false, "Retrieve a dad joke with an image")
	rootCmd.Flags().StringVarP(&imagePath, "imagePath", "p", "", "Path to save the image")

	rootCmd.Flags().BoolVarP(&slack, "asSlackMessage", "s", false, "Retrieve a dad joke as a slack message")

	rootCmd.MarkFlagsRequiredTogether("asImage", "imagePath")
}

func Execute() error {
	return rootCmd.Execute()
}
