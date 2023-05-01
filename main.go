package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/option"
)


func main() {
	ctx := context.Background()
	cloudbillingService, err := cloudbilling.NewService(ctx,option.WithCredentialsFile("serviceAccountKey.json"))
	if err != nil {
		log.Fatal("Failed to create cloudbilling service",err)
		return
	}
	request := cloudbillingService.Projects.GetBillingInfo("projects/readytoworkjapan")
	response, err := request.Do()
	if err != nil {
		log.Fatal("Failed to make a request",err)
		return
	}

	fmt.Println(request)
	resp,err := response.MarshalJSON()
	fmt.Printf("%+v\n",string(resp) )
}

// func getBillingInfo(projectID string, date string) (*cloudbilling.ProjectBillingInfo, error) {
//     ctx := context.Background()
//     billingService, err := cloudbilling.NewService(ctx, option.WithCredentialsFile("path/to/service-account.json"))
//     if err != nil {
//         return nil, fmt.Errorf("failed to create billing service: %v", err)
//     }

//     request := billingService.Projects.GetBillingInfo("projects/" + projectID)
//     request.BillingAccountName("billingAccounts/xxxxxx-xxxxxx-xxxxxx") // Replace with your billing account ID
//     request.CurrencyCode("USD") // Replace with your currency code
//     request.EffectiveTime(date + "T00:00:00Z")
//     response, err := request.Do()
//     if err != nil {
//         return nil, fmt.Errorf("failed to fetch billing info: %v", err)
//     }

//     return response, nil
// }

// func isCostExceeded(projectID string, date string, threshold float64) (bool, error) {
//     billingInfo, err := getBillingInfo(projectID, date)
//     if err != nil {
//         return false, err
//     }

//     totalCost := 0.0
//     for _, bucket := range billingInfo.BillingAccount.BillingEnabledProjects[0].ProjectBillingInfo.AllUpdates.BucketSummaries {
//         if bucket.BucketId == "cost" {
//             totalCost = bucket.CostAmount.Nanos / 1000000000.0
//         }
//     }

//     return totalCost > threshold, nil
// }
