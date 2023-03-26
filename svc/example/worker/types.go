package example_worker

// Predefined task types.
const (
	SendExampleTask = "send_{%example%}"
)

type (
	SendExamplePayload struct {
		ExampleID    string `json:"{%example%}_id"`
		Email        string `json:"email"`
		Role         string `json:"role"`
		MerchantName string `json:"merchant_name"`
	}
)
