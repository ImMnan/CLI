/*
Copyright © 2024 Manan Patel - github.com/immnan
*/
package get

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// integrationsCmd represents the integrations command
var integrationsCmd = &cobra.Command{
	Use:   "integrations",
	Short: "[>]Get list of integrations in a team",
	Long: `Use this command to list integrations for your Runscope/Api-Monitoring team. This is a Runscope only command. Ouputs integration Id, type, etc.
	
	For example: [bmgo get -t <team_uuid> integrations]
	For default: [bmgo get --tm integrations]`,
	Run: func(cmd *cobra.Command, args []string) {
		tm, _ := cmd.Flags().GetBool("tm")
		var teamId string
		if tm {
			teamId = defaultTeam()
		} else {
			teamId, _ = cmd.Flags().GetString("teamid")
		}
		rawOutput, _ := cmd.Flags().GetBool("raw")
		getIntegrationsTm(teamId, rawOutput)
	},
}

func init() {
	GetCmd.AddCommand(integrationsCmd)
}

type integrationsResponse struct {
	Data  []integrationsData `json:"data"`
	Error errorResult        `json:"error"`
}
type integrationsData struct {
	Uuid        string `json:"uuid"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

func getIntegrationsTm(teamId string, rawOutput bool) {
	Bearer := fmt.Sprintf("Bearer %v", GetPersonalAccessToken())
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.runscope.com/teams/"+teamId+"/integrations", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", Bearer)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if rawOutput {
		fmt.Printf("%s\n", bodyText)
	} else {
		var responseObjectIntegrations integrationsResponse
		json.Unmarshal(bodyText, &responseObjectIntegrations)
		if responseObjectIntegrations.Error.Code == 0 {
			//	fmt.Printf("\n%-40s %-15s %-10s\n", "UUID", "TYPE", "DESCRIPTION")
			for i := 0; i < len(responseObjectIntegrations.Data); i++ {
				integrationId := responseObjectIntegrations.Data[i].Uuid
				integrationsType := responseObjectIntegrations.Data[i].Type
				integrationsDesc := responseObjectIntegrations.Data[i].Description
				//	fmt.Printf("\n%-40s %-15s %-10s\n", integrationId, integrationsType, integrationsDesc)
				fmt.Printf("\nUUID: %s\nTYPE: %s\nDESCRIPTION: %s\n", integrationId, integrationsType, integrationsDesc)
			}
			fmt.Println("\n-")
		} else {
			errorCode := responseObjectIntegrations.Error.Status
			errorMessage := responseObjectIntegrations.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}

}
