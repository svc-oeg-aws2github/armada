package queue

import (
	"fmt"

	"github.com/G-Research/armada/pkg/api"
	"github.com/G-Research/armada/pkg/client"
	"github.com/spf13/cobra"
)

func Delete() *cobra.Command {
	command := &cobra.Command{
		Use:          "queue",
		Short:        "Delete existing queue",
		Long:         "Deletes queue if it exists, the queue needs to be empty at the time of deletion.",
		SilenceUsage: true,
	}

	command.Flags().SortFlags = false
	command.Flags().StringP("queueName", "n", "", "[required] Queue's name that will be deleted")
	command.MarkFlagRequired("queueName")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		queueName, err := cmd.Flags().GetString("queueName")
		if err != nil {
			return fmt.Errorf("failed to retrieve name value: %s", err)
		}

		apiConnectionDetails := client.ExtractCommandlineArmadaApiConnectionDetails()

		conn, err := client.CreateApiConnection(apiConnectionDetails)
		if err != nil {
			return fmt.Errorf("failed to connect to api because %s", err)
		}
		defer conn.Close()

		if err := client.DeleteQueue(api.NewSubmitClient(conn), queueName); err != nil {
			return fmt.Errorf("failed to delete queue because %s", err)
		}

		cmd.Printf("Queue %s deleted or did not exist.", queueName)
		return nil
	}

	return command
}
