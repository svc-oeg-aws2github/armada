package queue

import (
	"fmt"

	"github.com/G-Research/armada/pkg/api"
	"github.com/G-Research/armada/pkg/client"
	"github.com/spf13/cobra"
)

func Update() *cobra.Command {
	command := &cobra.Command{
		Use:   "queue",
		Short: "Update existing queue",
		Long:  "Update settings of existing queue",
	}

	command.Flags().SortFlags = false
	command.Flags().StringP("queueName", "n", "", "Queue name")
	command.MarkFlagRequired("queueName")

	command.Flags().Float64("priorityFactor", 1, "Set queue priority factor - lower number makes queue more important, must be > 0.")
	command.Flags().StringSlice("owners", []string{}, "Comma separated list of queue owners, defaults to current user.")
	command.Flags().StringSlice("groupOwners", []string{}, "Comma separated list of queue group owners, defaults to empty list.")
	command.Flags().StringToString("resourceLimits", map[string]string{},
		"Command separated list of resource limits pairs, defaults to empty list. Example: --resourceLimits cpu=0.3,memory=0.2",
	)

	command.RunE = func(cmd *cobra.Command, args []string) error {
		queueName, err := cmd.Flags().GetString("queueName")
		if err != nil {
			return fmt.Errorf("failed to retrieve name value: %s", err)
		}

		priority, err := cmd.Flags().GetFloat64("priorityFactor")
		if err != nil {
			return fmt.Errorf("failed to retrieve priorityFactor value: %s", err)
		}

		owners, err := cmd.Flags().GetStringSlice("owners")
		if err != nil {
			return fmt.Errorf("failed to retrieve owners value: %s", err)
		}

		groups, err := cmd.Flags().GetStringSlice("groupOwners")
		if err != nil {
			return fmt.Errorf("failed to retrieve groupOwners value: %s", err)
		}

		resourceLimits, err := FlagGetStringToString(cmd.Flags().GetStringToString).ToFloat64("resourceLimits")
		if err != nil {
			return fmt.Errorf("failed to retrieve resourceLimits value: %s", err)
		}

		apiConnectionDetails := client.ExtractCommandlineArmadaApiConnectionDetails()
		conn, err := client.CreateApiConnection(apiConnectionDetails)
		if err != nil {
			return fmt.Errorf("failed to connect to api because %s", err)
		}
		defer conn.Close()

		submissionClient := api.NewSubmitClient(conn)

		queue := &api.Queue{
			Name:           queueName,
			PriorityFactor: priority,
			UserOwners:     owners,
			GroupOwners:    groups,
			ResourceLimits: resourceLimits,
		}

		if err = client.UpdateQueue(submissionClient, queue); err != nil {
			return fmt.Errorf("failed to update queue with name %s. %s", queueName, err)
		}

		cmd.Printf("Queue %s updated", queue.Name)
		return nil
	}

	return command
}
