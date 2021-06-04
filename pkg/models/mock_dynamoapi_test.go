package models

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	DeleteItemMock func(*dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
	GetItemMock    func(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	PutItemMock    func(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	UpdateItemMock func(*dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
)

type DynamoServiceMock struct{}

func (d DynamoServiceMock) BatchExecuteStatement(*dynamodb.BatchExecuteStatementInput) (*dynamodb.BatchExecuteStatementOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) BatchExecuteStatementWithContext(aws.Context, *dynamodb.BatchExecuteStatementInput, ...request.Option) (*dynamodb.BatchExecuteStatementOutput, error) {
	return nil, nil
}
func (d DynamoServiceMock) BatchExecuteStatementRequest(*dynamodb.BatchExecuteStatementInput) (*request.Request, *dynamodb.BatchExecuteStatementOutput) {
	return nil, nil
}

func (d DynamoServiceMock) BatchGetItem(*dynamodb.BatchGetItemInput) (*dynamodb.BatchGetItemOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) BatchGetItemWithContext(aws.Context, *dynamodb.BatchGetItemInput, ...request.Option) (*dynamodb.BatchGetItemOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) BatchGetItemRequest(*dynamodb.BatchGetItemInput) (*request.Request, *dynamodb.BatchGetItemOutput) {
	return nil, nil
}

func (d DynamoServiceMock) BatchGetItemPages(*dynamodb.BatchGetItemInput, func(*dynamodb.BatchGetItemOutput, bool) bool) error {
	return errors.New("unimplemented")
}
func (d DynamoServiceMock) BatchGetItemPagesWithContext(aws.Context, *dynamodb.BatchGetItemInput, func(*dynamodb.BatchGetItemOutput, bool) bool, ...request.Option) error {
	return errors.New("unimplemented")
}

func (d DynamoServiceMock) BatchWriteItem(*dynamodb.BatchWriteItemInput) (*dynamodb.BatchWriteItemOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) BatchWriteItemWithContext(aws.Context, *dynamodb.BatchWriteItemInput, ...request.Option) (*dynamodb.BatchWriteItemOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) BatchWriteItemRequest(*dynamodb.BatchWriteItemInput) (*request.Request, *dynamodb.BatchWriteItemOutput) {
	return nil, nil
}

func (d DynamoServiceMock) CreateBackup(*dynamodb.CreateBackupInput) (*dynamodb.CreateBackupOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) CreateBackupWithContext(aws.Context, *dynamodb.CreateBackupInput, ...request.Option) (*dynamodb.CreateBackupOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) CreateBackupRequest(*dynamodb.CreateBackupInput) (*request.Request, *dynamodb.CreateBackupOutput) {
	return nil, nil
}

func (d DynamoServiceMock) CreateGlobalTable(*dynamodb.CreateGlobalTableInput) (*dynamodb.CreateGlobalTableOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) CreateGlobalTableWithContext(aws.Context, *dynamodb.CreateGlobalTableInput, ...request.Option) (*dynamodb.CreateGlobalTableOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) CreateGlobalTableRequest(*dynamodb.CreateGlobalTableInput) (*request.Request, *dynamodb.CreateGlobalTableOutput) {
	return nil, nil
}

func (d DynamoServiceMock) CreateTable(*dynamodb.CreateTableInput) (*dynamodb.CreateTableOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) CreateTableWithContext(aws.Context, *dynamodb.CreateTableInput, ...request.Option) (*dynamodb.CreateTableOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) CreateTableRequest(*dynamodb.CreateTableInput) (*request.Request, *dynamodb.CreateTableOutput) {
	return nil, nil
}

func (d DynamoServiceMock) DeleteBackup(*dynamodb.DeleteBackupInput) (*dynamodb.DeleteBackupOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DeleteBackupWithContext(aws.Context, *dynamodb.DeleteBackupInput, ...request.Option) (*dynamodb.DeleteBackupOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DeleteBackupRequest(*dynamodb.DeleteBackupInput) (*request.Request, *dynamodb.DeleteBackupOutput) {
	return nil, nil
}

func (d DynamoServiceMock) DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	return DeleteItemMock(input)
}
func (d DynamoServiceMock) DeleteItemWithContext(aws.Context, *dynamodb.DeleteItemInput, ...request.Option) (*dynamodb.DeleteItemOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DeleteItemRequest(*dynamodb.DeleteItemInput) (*request.Request, *dynamodb.DeleteItemOutput) {
	return nil, nil
}

func (d DynamoServiceMock) DeleteTable(*dynamodb.DeleteTableInput) (*dynamodb.DeleteTableOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DeleteTableWithContext(aws.Context, *dynamodb.DeleteTableInput, ...request.Option) (*dynamodb.DeleteTableOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DeleteTableRequest(*dynamodb.DeleteTableInput) (*request.Request, *dynamodb.DeleteTableOutput) {
	return nil, nil
}

func (d DynamoServiceMock) DescribeBackup(*dynamodb.DescribeBackupInput) (*dynamodb.DescribeBackupOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DescribeBackupWithContext(aws.Context, *dynamodb.DescribeBackupInput, ...request.Option) (*dynamodb.DescribeBackupOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DescribeBackupRequest(*dynamodb.DescribeBackupInput) (*request.Request, *dynamodb.DescribeBackupOutput) {
	return nil, nil
}

func (d DynamoServiceMock) DescribeContinuousBackups(*dynamodb.DescribeContinuousBackupsInput) (*dynamodb.DescribeContinuousBackupsOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DescribeContinuousBackupsWithContext(aws.Context, *dynamodb.DescribeContinuousBackupsInput, ...request.Option) (*dynamodb.DescribeContinuousBackupsOutput, error) {
	return nil, nil
}
func (d DynamoServiceMock) DescribeContinuousBackupsRequest(*dynamodb.DescribeContinuousBackupsInput) (*request.Request, *dynamodb.DescribeContinuousBackupsOutput) {
	return nil, nil
}

func (d DynamoServiceMock) DescribeContributorInsights(*dynamodb.DescribeContributorInsightsInput) (*dynamodb.DescribeContributorInsightsOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DescribeContributorInsightsWithContext(aws.Context, *dynamodb.DescribeContributorInsightsInput, ...request.Option) (*dynamodb.DescribeContributorInsightsOutput, error) {
	return nil, nil
}
func (d DynamoServiceMock) DescribeContributorInsightsRequest(*dynamodb.DescribeContributorInsightsInput) (*request.Request, *dynamodb.DescribeContributorInsightsOutput) {
	return nil, nil
}

func (d DynamoServiceMock) DescribeEndpoints(*dynamodb.DescribeEndpointsInput) (*dynamodb.DescribeEndpointsOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DescribeEndpointsWithContext(aws.Context, *dynamodb.DescribeEndpointsInput, ...request.Option) (*dynamodb.DescribeEndpointsOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DescribeEndpointsRequest(*dynamodb.DescribeEndpointsInput) (*request.Request, *dynamodb.DescribeEndpointsOutput) {
	return nil, nil
}

func (d DynamoServiceMock) DescribeExport(*dynamodb.DescribeExportInput) (*dynamodb.DescribeExportOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DescribeExportWithContext(aws.Context, *dynamodb.DescribeExportInput, ...request.Option) (*dynamodb.DescribeExportOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DescribeExportRequest(*dynamodb.DescribeExportInput) (*request.Request, *dynamodb.DescribeExportOutput) {
	return nil, nil
}

func (d DynamoServiceMock) DescribeGlobalTable(*dynamodb.DescribeGlobalTableInput) (*dynamodb.DescribeGlobalTableOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DescribeGlobalTableWithContext(aws.Context, *dynamodb.DescribeGlobalTableInput, ...request.Option) (*dynamodb.DescribeGlobalTableOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DescribeGlobalTableRequest(*dynamodb.DescribeGlobalTableInput) (*request.Request, *dynamodb.DescribeGlobalTableOutput) {
	return nil, nil
}

func (d DynamoServiceMock) DescribeGlobalTableSettings(*dynamodb.DescribeGlobalTableSettingsInput) (*dynamodb.DescribeGlobalTableSettingsOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DescribeGlobalTableSettingsWithContext(aws.Context, *dynamodb.DescribeGlobalTableSettingsInput, ...request.Option) (*dynamodb.DescribeGlobalTableSettingsOutput, error) {
	return nil, nil
}
func (d DynamoServiceMock) DescribeGlobalTableSettingsRequest(*dynamodb.DescribeGlobalTableSettingsInput) (*request.Request, *dynamodb.DescribeGlobalTableSettingsOutput) {
	return nil, nil
}

func (d DynamoServiceMock) DescribeKinesisStreamingDestination(*dynamodb.DescribeKinesisStreamingDestinationInput) (*dynamodb.DescribeKinesisStreamingDestinationOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DescribeKinesisStreamingDestinationWithContext(aws.Context, *dynamodb.DescribeKinesisStreamingDestinationInput, ...request.Option) (*dynamodb.DescribeKinesisStreamingDestinationOutput, error) {
	return nil, nil
}
func (d DynamoServiceMock) DescribeKinesisStreamingDestinationRequest(*dynamodb.DescribeKinesisStreamingDestinationInput) (*request.Request, *dynamodb.DescribeKinesisStreamingDestinationOutput) {
	return nil, nil
}

func (d DynamoServiceMock) DescribeLimits(*dynamodb.DescribeLimitsInput) (*dynamodb.DescribeLimitsOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DescribeLimitsWithContext(aws.Context, *dynamodb.DescribeLimitsInput, ...request.Option) (*dynamodb.DescribeLimitsOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DescribeLimitsRequest(*dynamodb.DescribeLimitsInput) (*request.Request, *dynamodb.DescribeLimitsOutput) {
	return nil, nil
}

func (d DynamoServiceMock) DescribeTable(*dynamodb.DescribeTableInput) (*dynamodb.DescribeTableOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DescribeTableWithContext(aws.Context, *dynamodb.DescribeTableInput, ...request.Option) (*dynamodb.DescribeTableOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DescribeTableRequest(*dynamodb.DescribeTableInput) (*request.Request, *dynamodb.DescribeTableOutput) {
	return nil, nil
}

func (d DynamoServiceMock) DescribeTableReplicaAutoScaling(*dynamodb.DescribeTableReplicaAutoScalingInput) (*dynamodb.DescribeTableReplicaAutoScalingOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DescribeTableReplicaAutoScalingWithContext(aws.Context, *dynamodb.DescribeTableReplicaAutoScalingInput, ...request.Option) (*dynamodb.DescribeTableReplicaAutoScalingOutput, error) {
	return nil, nil
}
func (d DynamoServiceMock) DescribeTableReplicaAutoScalingRequest(*dynamodb.DescribeTableReplicaAutoScalingInput) (*request.Request, *dynamodb.DescribeTableReplicaAutoScalingOutput) {
	return nil, nil
}

func (d DynamoServiceMock) DescribeTimeToLive(*dynamodb.DescribeTimeToLiveInput) (*dynamodb.DescribeTimeToLiveOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DescribeTimeToLiveWithContext(aws.Context, *dynamodb.DescribeTimeToLiveInput, ...request.Option) (*dynamodb.DescribeTimeToLiveOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DescribeTimeToLiveRequest(*dynamodb.DescribeTimeToLiveInput) (*request.Request, *dynamodb.DescribeTimeToLiveOutput) {
	return nil, nil
}

func (d DynamoServiceMock) DisableKinesisStreamingDestination(*dynamodb.DisableKinesisStreamingDestinationInput) (*dynamodb.DisableKinesisStreamingDestinationOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) DisableKinesisStreamingDestinationWithContext(aws.Context, *dynamodb.DisableKinesisStreamingDestinationInput, ...request.Option) (*dynamodb.DisableKinesisStreamingDestinationOutput, error) {
	return nil, nil
}
func (d DynamoServiceMock) DisableKinesisStreamingDestinationRequest(*dynamodb.DisableKinesisStreamingDestinationInput) (*request.Request, *dynamodb.DisableKinesisStreamingDestinationOutput) {
	return nil, nil
}

func (d DynamoServiceMock) EnableKinesisStreamingDestination(*dynamodb.EnableKinesisStreamingDestinationInput) (*dynamodb.EnableKinesisStreamingDestinationOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) EnableKinesisStreamingDestinationWithContext(aws.Context, *dynamodb.EnableKinesisStreamingDestinationInput, ...request.Option) (*dynamodb.EnableKinesisStreamingDestinationOutput, error) {
	return nil, nil
}
func (d DynamoServiceMock) EnableKinesisStreamingDestinationRequest(*dynamodb.EnableKinesisStreamingDestinationInput) (*request.Request, *dynamodb.EnableKinesisStreamingDestinationOutput) {
	return nil, nil
}

func (d DynamoServiceMock) ExecuteStatement(*dynamodb.ExecuteStatementInput) (*dynamodb.ExecuteStatementOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) ExecuteStatementWithContext(aws.Context, *dynamodb.ExecuteStatementInput, ...request.Option) (*dynamodb.ExecuteStatementOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) ExecuteStatementRequest(*dynamodb.ExecuteStatementInput) (*request.Request, *dynamodb.ExecuteStatementOutput) {
	return nil, nil
}

func (d DynamoServiceMock) ExecuteTransaction(*dynamodb.ExecuteTransactionInput) (*dynamodb.ExecuteTransactionOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) ExecuteTransactionWithContext(aws.Context, *dynamodb.ExecuteTransactionInput, ...request.Option) (*dynamodb.ExecuteTransactionOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) ExecuteTransactionRequest(*dynamodb.ExecuteTransactionInput) (*request.Request, *dynamodb.ExecuteTransactionOutput) {
	return nil, nil
}

func (d DynamoServiceMock) ExportTableToPointInTime(*dynamodb.ExportTableToPointInTimeInput) (*dynamodb.ExportTableToPointInTimeOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) ExportTableToPointInTimeWithContext(aws.Context, *dynamodb.ExportTableToPointInTimeInput, ...request.Option) (*dynamodb.ExportTableToPointInTimeOutput, error) {
	return nil, nil
}
func (d DynamoServiceMock) ExportTableToPointInTimeRequest(*dynamodb.ExportTableToPointInTimeInput) (*request.Request, *dynamodb.ExportTableToPointInTimeOutput) {
	return nil, nil
}

func (d DynamoServiceMock) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return GetItemMock(input)
}
func (d DynamoServiceMock) GetItemWithContext(aws.Context, *dynamodb.GetItemInput, ...request.Option) (*dynamodb.GetItemOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) GetItemRequest(*dynamodb.GetItemInput) (*request.Request, *dynamodb.GetItemOutput) {
	return nil, nil
}

func (d DynamoServiceMock) ListBackups(*dynamodb.ListBackupsInput) (*dynamodb.ListBackupsOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) ListBackupsWithContext(aws.Context, *dynamodb.ListBackupsInput, ...request.Option) (*dynamodb.ListBackupsOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) ListBackupsRequest(*dynamodb.ListBackupsInput) (*request.Request, *dynamodb.ListBackupsOutput) {
	return nil, nil
}

func (d DynamoServiceMock) ListContributorInsights(*dynamodb.ListContributorInsightsInput) (*dynamodb.ListContributorInsightsOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) ListContributorInsightsWithContext(aws.Context, *dynamodb.ListContributorInsightsInput, ...request.Option) (*dynamodb.ListContributorInsightsOutput, error) {
	return nil, nil
}
func (d DynamoServiceMock) ListContributorInsightsRequest(*dynamodb.ListContributorInsightsInput) (*request.Request, *dynamodb.ListContributorInsightsOutput) {
	return nil, nil
}

func (d DynamoServiceMock) ListContributorInsightsPages(*dynamodb.ListContributorInsightsInput, func(*dynamodb.ListContributorInsightsOutput, bool) bool) error {
	return errors.New("unimplemented")
}
func (d DynamoServiceMock) ListContributorInsightsPagesWithContext(aws.Context, *dynamodb.ListContributorInsightsInput, func(*dynamodb.ListContributorInsightsOutput, bool) bool, ...request.Option) error {
	return errors.New("unimplemented")
}

func (d DynamoServiceMock) ListExports(*dynamodb.ListExportsInput) (*dynamodb.ListExportsOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) ListExportsWithContext(aws.Context, *dynamodb.ListExportsInput, ...request.Option) (*dynamodb.ListExportsOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) ListExportsRequest(*dynamodb.ListExportsInput) (*request.Request, *dynamodb.ListExportsOutput) {
	return nil, nil
}

func (d DynamoServiceMock) ListExportsPages(*dynamodb.ListExportsInput, func(*dynamodb.ListExportsOutput, bool) bool) error {
	return errors.New("unimplemented")
}
func (d DynamoServiceMock) ListExportsPagesWithContext(aws.Context, *dynamodb.ListExportsInput, func(*dynamodb.ListExportsOutput, bool) bool, ...request.Option) error {
	return errors.New("unimplemented")
}

func (d DynamoServiceMock) ListGlobalTables(*dynamodb.ListGlobalTablesInput) (*dynamodb.ListGlobalTablesOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) ListGlobalTablesWithContext(aws.Context, *dynamodb.ListGlobalTablesInput, ...request.Option) (*dynamodb.ListGlobalTablesOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) ListGlobalTablesRequest(*dynamodb.ListGlobalTablesInput) (*request.Request, *dynamodb.ListGlobalTablesOutput) {
	return nil, nil
}

func (d DynamoServiceMock) ListTables(*dynamodb.ListTablesInput) (*dynamodb.ListTablesOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) ListTablesWithContext(aws.Context, *dynamodb.ListTablesInput, ...request.Option) (*dynamodb.ListTablesOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) ListTablesRequest(*dynamodb.ListTablesInput) (*request.Request, *dynamodb.ListTablesOutput) {
	return nil, nil
}

func (d DynamoServiceMock) ListTablesPages(*dynamodb.ListTablesInput, func(*dynamodb.ListTablesOutput, bool) bool) error {
	return errors.New("unimplemented")
}
func (d DynamoServiceMock) ListTablesPagesWithContext(aws.Context, *dynamodb.ListTablesInput, func(*dynamodb.ListTablesOutput, bool) bool, ...request.Option) error {
	return errors.New("unimplemented")
}

func (d DynamoServiceMock) ListTagsOfResource(*dynamodb.ListTagsOfResourceInput) (*dynamodb.ListTagsOfResourceOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) ListTagsOfResourceWithContext(aws.Context, *dynamodb.ListTagsOfResourceInput, ...request.Option) (*dynamodb.ListTagsOfResourceOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) ListTagsOfResourceRequest(*dynamodb.ListTagsOfResourceInput) (*request.Request, *dynamodb.ListTagsOfResourceOutput) {
	return nil, nil
}

func (d DynamoServiceMock) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return PutItemMock(input)
}
func (d DynamoServiceMock) PutItemWithContext(aws.Context, *dynamodb.PutItemInput, ...request.Option) (*dynamodb.PutItemOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) PutItemRequest(*dynamodb.PutItemInput) (*request.Request, *dynamodb.PutItemOutput) {
	return nil, nil
}

func (d DynamoServiceMock) Query(*dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) QueryWithContext(aws.Context, *dynamodb.QueryInput, ...request.Option) (*dynamodb.QueryOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) QueryRequest(*dynamodb.QueryInput) (*request.Request, *dynamodb.QueryOutput) {
	return nil, nil
}

func (d DynamoServiceMock) QueryPages(*dynamodb.QueryInput, func(*dynamodb.QueryOutput, bool) bool) error {
	return errors.New("unimplemented")
}
func (d DynamoServiceMock) QueryPagesWithContext(aws.Context, *dynamodb.QueryInput, func(*dynamodb.QueryOutput, bool) bool, ...request.Option) error {
	return errors.New("unimplemented")
}

func (d DynamoServiceMock) RestoreTableFromBackup(*dynamodb.RestoreTableFromBackupInput) (*dynamodb.RestoreTableFromBackupOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) RestoreTableFromBackupWithContext(aws.Context, *dynamodb.RestoreTableFromBackupInput, ...request.Option) (*dynamodb.RestoreTableFromBackupOutput, error) {
	return nil, nil
}
func (d DynamoServiceMock) RestoreTableFromBackupRequest(*dynamodb.RestoreTableFromBackupInput) (*request.Request, *dynamodb.RestoreTableFromBackupOutput) {
	return nil, nil
}

func (d DynamoServiceMock) RestoreTableToPointInTime(*dynamodb.RestoreTableToPointInTimeInput) (*dynamodb.RestoreTableToPointInTimeOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) RestoreTableToPointInTimeWithContext(aws.Context, *dynamodb.RestoreTableToPointInTimeInput, ...request.Option) (*dynamodb.RestoreTableToPointInTimeOutput, error) {
	return nil, nil
}
func (d DynamoServiceMock) RestoreTableToPointInTimeRequest(*dynamodb.RestoreTableToPointInTimeInput) (*request.Request, *dynamodb.RestoreTableToPointInTimeOutput) {
	return nil, nil
}

func (d DynamoServiceMock) Scan(*dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) ScanWithContext(aws.Context, *dynamodb.ScanInput, ...request.Option) (*dynamodb.ScanOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) ScanRequest(*dynamodb.ScanInput) (*request.Request, *dynamodb.ScanOutput) {
	return nil, nil
}

func (d DynamoServiceMock) ScanPages(*dynamodb.ScanInput, func(*dynamodb.ScanOutput, bool) bool) error {
	return errors.New("unimplemented")
}
func (d DynamoServiceMock) ScanPagesWithContext(aws.Context, *dynamodb.ScanInput, func(*dynamodb.ScanOutput, bool) bool, ...request.Option) error {
	return errors.New("unimplemented")
}

func (d DynamoServiceMock) TagResource(*dynamodb.TagResourceInput) (*dynamodb.TagResourceOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) TagResourceWithContext(aws.Context, *dynamodb.TagResourceInput, ...request.Option) (*dynamodb.TagResourceOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) TagResourceRequest(*dynamodb.TagResourceInput) (*request.Request, *dynamodb.TagResourceOutput) {
	return nil, nil
}

func (d DynamoServiceMock) TransactGetItems(*dynamodb.TransactGetItemsInput) (*dynamodb.TransactGetItemsOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) TransactGetItemsWithContext(aws.Context, *dynamodb.TransactGetItemsInput, ...request.Option) (*dynamodb.TransactGetItemsOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) TransactGetItemsRequest(*dynamodb.TransactGetItemsInput) (*request.Request, *dynamodb.TransactGetItemsOutput) {
	return nil, nil
}

func (d DynamoServiceMock) TransactWriteItems(*dynamodb.TransactWriteItemsInput) (*dynamodb.TransactWriteItemsOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) TransactWriteItemsWithContext(aws.Context, *dynamodb.TransactWriteItemsInput, ...request.Option) (*dynamodb.TransactWriteItemsOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) TransactWriteItemsRequest(*dynamodb.TransactWriteItemsInput) (*request.Request, *dynamodb.TransactWriteItemsOutput) {
	return nil, nil
}

func (d DynamoServiceMock) UntagResource(*dynamodb.UntagResourceInput) (*dynamodb.UntagResourceOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) UntagResourceWithContext(aws.Context, *dynamodb.UntagResourceInput, ...request.Option) (*dynamodb.UntagResourceOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) UntagResourceRequest(*dynamodb.UntagResourceInput) (*request.Request, *dynamodb.UntagResourceOutput) {
	return nil, nil
}

func (d DynamoServiceMock) UpdateContinuousBackups(*dynamodb.UpdateContinuousBackupsInput) (*dynamodb.UpdateContinuousBackupsOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) UpdateContinuousBackupsWithContext(aws.Context, *dynamodb.UpdateContinuousBackupsInput, ...request.Option) (*dynamodb.UpdateContinuousBackupsOutput, error) {
	return nil, nil
}
func (d DynamoServiceMock) UpdateContinuousBackupsRequest(*dynamodb.UpdateContinuousBackupsInput) (*request.Request, *dynamodb.UpdateContinuousBackupsOutput) {
	return nil, nil
}

func (d DynamoServiceMock) UpdateContributorInsights(*dynamodb.UpdateContributorInsightsInput) (*dynamodb.UpdateContributorInsightsOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) UpdateContributorInsightsWithContext(aws.Context, *dynamodb.UpdateContributorInsightsInput, ...request.Option) (*dynamodb.UpdateContributorInsightsOutput, error) {
	return nil, nil
}
func (d DynamoServiceMock) UpdateContributorInsightsRequest(*dynamodb.UpdateContributorInsightsInput) (*request.Request, *dynamodb.UpdateContributorInsightsOutput) {
	return nil, nil
}

func (d DynamoServiceMock) UpdateGlobalTable(*dynamodb.UpdateGlobalTableInput) (*dynamodb.UpdateGlobalTableOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) UpdateGlobalTableWithContext(aws.Context, *dynamodb.UpdateGlobalTableInput, ...request.Option) (*dynamodb.UpdateGlobalTableOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) UpdateGlobalTableRequest(*dynamodb.UpdateGlobalTableInput) (*request.Request, *dynamodb.UpdateGlobalTableOutput) {
	return nil, nil
}

func (d DynamoServiceMock) UpdateGlobalTableSettings(*dynamodb.UpdateGlobalTableSettingsInput) (*dynamodb.UpdateGlobalTableSettingsOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) UpdateGlobalTableSettingsWithContext(aws.Context, *dynamodb.UpdateGlobalTableSettingsInput, ...request.Option) (*dynamodb.UpdateGlobalTableSettingsOutput, error) {
	return nil, nil
}
func (d DynamoServiceMock) UpdateGlobalTableSettingsRequest(*dynamodb.UpdateGlobalTableSettingsInput) (*request.Request, *dynamodb.UpdateGlobalTableSettingsOutput) {
	return nil, nil
}

func (d DynamoServiceMock) UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	return UpdateItemMock(input)
}
func (d DynamoServiceMock) UpdateItemWithContext(aws.Context, *dynamodb.UpdateItemInput, ...request.Option) (*dynamodb.UpdateItemOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) UpdateItemRequest(*dynamodb.UpdateItemInput) (*request.Request, *dynamodb.UpdateItemOutput) {
	return nil, nil
}

func (d DynamoServiceMock) UpdateTable(*dynamodb.UpdateTableInput) (*dynamodb.UpdateTableOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) UpdateTableWithContext(aws.Context, *dynamodb.UpdateTableInput, ...request.Option) (*dynamodb.UpdateTableOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) UpdateTableRequest(*dynamodb.UpdateTableInput) (*request.Request, *dynamodb.UpdateTableOutput) {
	return nil, nil
}

func (d DynamoServiceMock) UpdateTableReplicaAutoScaling(*dynamodb.UpdateTableReplicaAutoScalingInput) (*dynamodb.UpdateTableReplicaAutoScalingOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) UpdateTableReplicaAutoScalingWithContext(aws.Context, *dynamodb.UpdateTableReplicaAutoScalingInput, ...request.Option) (*dynamodb.UpdateTableReplicaAutoScalingOutput, error) {
	return nil, nil
}
func (d DynamoServiceMock) UpdateTableReplicaAutoScalingRequest(*dynamodb.UpdateTableReplicaAutoScalingInput) (*request.Request, *dynamodb.UpdateTableReplicaAutoScalingOutput) {
	return nil, nil
}

func (d DynamoServiceMock) UpdateTimeToLive(*dynamodb.UpdateTimeToLiveInput) (*dynamodb.UpdateTimeToLiveOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) UpdateTimeToLiveWithContext(aws.Context, *dynamodb.UpdateTimeToLiveInput, ...request.Option) (*dynamodb.UpdateTimeToLiveOutput, error) {
	return nil, errors.New("unimplemented")
}
func (d DynamoServiceMock) UpdateTimeToLiveRequest(*dynamodb.UpdateTimeToLiveInput) (*request.Request, *dynamodb.UpdateTimeToLiveOutput) {
	return nil, nil
}

func (d DynamoServiceMock) WaitUntilTableExists(*dynamodb.DescribeTableInput) error {
	return errors.New("unimplemented")
}
func (d DynamoServiceMock) WaitUntilTableExistsWithContext(aws.Context, *dynamodb.DescribeTableInput, ...request.WaiterOption) error {
	return errors.New("unimplemented")
}

func (d DynamoServiceMock) WaitUntilTableNotExists(*dynamodb.DescribeTableInput) error {
	return errors.New("unimplemented")
}
func (d DynamoServiceMock) WaitUntilTableNotExistsWithContext(aws.Context, *dynamodb.DescribeTableInput, ...request.WaiterOption) error {
	return errors.New("unimplemented")
}
