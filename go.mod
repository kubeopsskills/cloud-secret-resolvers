module github.com/kubeopsskills/cloud-secret-resolvers

go 1.15

require github.com/aws/aws-sdk-go v1.37.1

replace github.com/kubeopsskills/cloud-secret-resolvers/internal/csr => /internal/csr
replace github.com/kubeopsskills/cloud-secret-resolvers/internal/pkg => /internal/pkg
replace github.com/kubeopsskills/cloud-secret-resolvers/internal/pkg/provider => /internal/pkg/provider