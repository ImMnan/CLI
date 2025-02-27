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
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

// agentsCmd represents the agents command
var agentsCmd = &cobra.Command{
	Use:   "agents",
	Short: "Get agents within a private location",
	Long: `The command returns a list of created agents, you will need to provide a workspace id or a harborId to run the command. any server on which you install our agent is an agent within a Private location. These are your load generators. Formerly known as a 'ship'. The command returns a list of agents within a workspace or within a harborId id specified.  Outputs "SHIP ID", "STATE", etc.
	
	For example: [bmgo get -w <workspace id> agents --hid <harbour_id>] OR 
	             [bmgo get -w <workspace_id> agents]
				 [bmgo get -t <team_id> agents]
	For default: [bmgo get --ws agents]
	             [bmgo get --ws agents --hid <harbour id>]
				 [bmgo get --tm agents]`,
	Run: func(cmd *cobra.Command, args []string) {
		ws, _ := cmd.Flags().GetBool("ws")
		tm, _ := cmd.Flags().GetBool("tm")
		var workspaceId int
		var teamId string
		if ws {
			workspaceId = defaultWorkspace()
		} else {
			workspaceId, _ = cmd.Flags().GetInt("workspaceid")
		}
		if tm {
			teamId = defaultTeam()
		} else {
			teamId, _ = cmd.Flags().GetString("teamid")
		}
		rawOutput, _ := cmd.Flags().GetBool("raw")
		harbourId, _ := cmd.Flags().GetString("hid")
		switch {
		case workspaceId == 0 && harbourId == "" && teamId != "":
			getAgentsTm(teamId, rawOutput)
		case workspaceId != 0 && harbourId == "":
			getAgentsWs(workspaceId, rawOutput)
		case workspaceId != 0 && harbourId != "":
			getAgentsOpl(workspaceId, harbourId, rawOutput)
		default:
			cmd.Help()
		}
	},
}

func init() {
	GetCmd.AddCommand(agentsCmd)
	agentsCmd.Flags().String("hid", "", "Provide the harbour id")
}

// This function is because the API response for listing agents using harbour id is a struct & not a list/array to iterate over.
func getAgentsOpl(workspaceId int, harbourId string, rawOutput bool) {
	apiId, apiSecret := Getapikeys()
	workspaceIdStr := strconv.Itoa(workspaceId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/private-locations?workspaceId="+workspaceIdStr+"&limit=0", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(apiId, apiSecret)
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
		var responseBodyWsAgents oplsResponse
		json.Unmarshal(bodyText, &responseBodyWsAgents)
		if responseBodyWsAgents.Error.Code == 0 {
			for i := 0; i < len(responseBodyWsAgents.Result); i++ {
				oplId := responseBodyWsAgents.Result[i].Id
				if oplId == harbourId {
					fmt.Printf("For OPL/HARBOUR %v & NAMED %v:\n", oplId, responseBodyWsAgents.Result[i].Name)
					tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
					// Print headers
					fmt.Fprintln(tabWriter, "SHIP ID\tSTATE\tLAST BEAT\tNAME")
					//fmt.Printf("\n%-28s %-8s %-18s %-10s\n", "SHIP ID", "STATE", "LAST BEAT", "NAME")
					for f := 0; f < len(responseBodyWsAgents.Result[i].Ships); f++ {
						shipId := responseBodyWsAgents.Result[i].Ships[f].Id
						shipName := responseBodyWsAgents.Result[i].Ships[f].Name
						shipStatus := responseBodyWsAgents.Result[i].Ships[f].State
						shipLastHeartBeatEp := int64(responseBodyWsAgents.Result[i].Ships[f].LastHeartBeat)
						//	shipLastHeartBeat := time.Unix(shipLastHeartBeatEp, 0)
						if shipLastHeartBeatEp != 0 {
							shipLastHeartBeatStr := fmt.Sprint(time.Unix(shipLastHeartBeatEp, 0))
							//	fmt.Printf("\n%-28s %-8s %-18s %-10s", shipId, shipStatus, shipLastHeartBeatStr[0:16], shipName)
							fmt.Fprintf(tabWriter, "%s\t%s\t%s\t%s\n", shipId, shipStatus, shipLastHeartBeatStr[0:16], shipName)
						} else {
							shipLastHeartBeat := shipLastHeartBeatEp
							//	fmt.Printf("\n%-28s %-8s %-18d %-10s", shipId, shipStatus, shipLastHeartBeat, shipName)
							fmt.Fprintf(tabWriter, "%s\t%s\t%d\t%s\n", shipId, shipStatus, shipLastHeartBeat, shipName)
						}
					}
					tabWriter.Flush()
					fmt.Println("-")
				} else {
					continue
				}
			}
		} else {
			errorCode := responseBodyWsAgents.Error.Code
			errorMessage := responseBodyWsAgents.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}
func getAgentsWs(workspaceId int, rawOutput bool) {
	apiId, apiSecret := Getapikeys()
	workspaceIdStr := strconv.Itoa(workspaceId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/private-locations?workspaceId="+workspaceIdStr+"&limit=0", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(apiId, apiSecret)
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
		var responseBodyWsAgents oplsResponse
		json.Unmarshal(bodyText, &responseBodyWsAgents)
		if responseBodyWsAgents.Error.Code == 0 {
			for i := 0; i < len(responseBodyWsAgents.Result); i++ {
				fmt.Printf("For OPL/HARBOUR %v & NAMED %v:\n", responseBodyWsAgents.Result[i].Id, responseBodyWsAgents.Result[i].Name)
				tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				// Print headers
				fmt.Fprintln(tabWriter, "SHIP ID\tSTATE\tLAST BEAT\tNAME")
				for f := 0; f < len(responseBodyWsAgents.Result[i].Ships); f++ {
					shipId := responseBodyWsAgents.Result[i].Ships[f].Id
					shipName := responseBodyWsAgents.Result[i].Ships[f].Name
					shipStatus := responseBodyWsAgents.Result[i].Ships[f].State
					shipLastHeartBeatEp := int64(responseBodyWsAgents.Result[i].Ships[f].LastHeartBeat)
					//	shipLastHeartBeat := time.Unix(shipLastHeartBeatEp, 0)
					if shipLastHeartBeatEp != 0 {
						shipLastHeartBeatStr := fmt.Sprint(time.Unix(shipLastHeartBeatEp, 0))
						fmt.Fprintf(tabWriter, "%s\t%s\t%s\t%s\n", shipId, shipStatus, shipLastHeartBeatStr[0:16], shipName)
					} else {
						shipLastHeartBeat := shipLastHeartBeatEp
						fmt.Fprintf(tabWriter, "%s\t%s\t%d\t%s\n", shipId, shipStatus, shipLastHeartBeat, shipName)
					}
				}
				tabWriter.Flush()
				fmt.Println("\n-")
			}
		} else {
			errorCode := responseBodyWsAgents.Error.Code
			errorMessage := responseBodyWsAgents.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}

type rsAgentResponse struct {
	Meta meta          `json:"meta"`
	Data []resultsData `json:"data"`
}
type meta struct {
	Status string `json:"status"`
}
type resultsData struct {
	Agent_id string `json:"agent_id"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	HostOs   string `json:"host_os"`
}

func getAgentsTm(teamId string, rawOutput bool) {
	Bearer := fmt.Sprintf("Bearer %v", GetPersonalAccessToken())
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.runscope.com/v1/teams/"+teamId+"/agents", nil)
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
		var responseBodyTmAgents rsAgentResponse
		json.Unmarshal(bodyText, &responseBodyTmAgents)
		if responseBodyTmAgents.Meta.Status == "success" {
			fmt.Printf("\n%-37s %-6s %-22s %-5s\n", "AGENT ID", "OS", "VERSION", "NAME")
			tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			// Print headers
			fmt.Fprintln(tabWriter, "AGENT ID\tOS\tVERSION\tNAME")
			for i := 0; i < len(responseBodyTmAgents.Data); i++ {
				agentIdtm := responseBodyTmAgents.Data[i].Agent_id
				agentNametm := responseBodyTmAgents.Data[i].Name
				agentVersiontm := responseBodyTmAgents.Data[i].Version
				agentHostOstm := responseBodyTmAgents.Data[i].HostOs
				//			fmt.Printf("\n%-37s %-6s %-22s %-5s", agentIdtm, agentHostOstm, agentVersiontm, agentNametm)
				fmt.Fprintf(tabWriter, "%s\t%s\t%s\t%s\n", agentIdtm, agentHostOstm, agentVersiontm, agentNametm)
			}
			tabWriter.Flush()
			fmt.Println("\n-")
		} else {
			fmt.Println(responseBodyTmAgents.Meta.Status)
		}
	}
}
