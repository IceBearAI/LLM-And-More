package fastchat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/http"
	"net/http/httputil"
	"testing"
)

func initSvc() Service {
	return New(Config{
		OpenAiEndpoint:  "https://api.openai.com/v1",
		OpenAiToken:     "",
		OpenAiModel:     "gpt-3.5-turbo",
		LocalAiEndpoint: "http://fastchat:8000",
		LocalAiToken:    "sk-001",
		OpenAiOrgId:     "",
		Debug:           true,
	}, []kithttp.ClientOption{
		kithttp.ClientBefore(func(ctx context.Context, request *http.Request) context.Context {
			dumpRequest, err := httputil.DumpRequest(request, true)
			if err != nil {
				return ctx
			}
			fmt.Println(string(dumpRequest))
			return ctx
		}),
		kithttp.ClientAfter(func(ctx context.Context, response *http.Response) context.Context {
			dumpResponse, err := httputil.DumpResponse(response, true)
			if err != nil {
				return ctx
			}
			fmt.Println(string(dumpResponse))
			return ctx
		}),
	},
	)
}

func TestService_ChatProcess(t *testing.T) {
	ctx := context.Background()
	svc := initSvc()
	process, _, err := svc.ChatCompletion(ctx, "qwen-72b-chat-int4", []openai.ChatCompletionMessage{
		{
			Role:    "user",
			Content: "给你取个名字，你就叫AIGC吧！",
		},
		{
			Role:    "assistant",
			Content: "好的，我现在可以被称为智语。谢谢您为我取这个名字！",
		}, {
			Role:    "user",
			Content: "你叫什么名字？",
		},
	}, 0.95, 0.7, 0.5, 0.5, 1024, 1, nil, "", nil, nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(process)
}

func TestService_ChatCompletionStream(t *testing.T) {
	ctx := context.Background()
	svc := initSvc()
	ctx = context.WithValue(ctx, ContextKeyPlatform, PlatformOpenAI)
	stream, _, err := svc.ChatCompletionStream(ctx, "gpt-3.5-turbo", []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "你就是训练的大型语言模型AIGC。仔细按照用户的说明进行操作。使用Markdown响应。",
		}, {
			Role:    openai.ChatMessageRoleUser,
			Content: "你叫什么名字？",
		},
	}, 0.95, 0.7, 0, 0, 1024, 1, nil, "", nil, nil)
	if err != nil {
		t.Error(err)
		return
	}
	defer stream.Close()
	for {
		recv, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			t.Log(err)
			break
		}
		if errors.Is(err, openai.ErrTooManyEmptyStreamMessages) {
			t.Error(err)
			break
		}
		if err != nil {
			t.Error(err)
			break
		}
		t.Log(recv.Choices[0].Delta.Content)
	}
	t.Log("ok")
}

func TestService_Models(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ContextKeyPlatform, PlatformOpenAI)
	svc := initSvc()
	models, err := svc.Models(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range models {
		b, _ := json.Marshal(v)
		t.Log(string(b))
	}
}

func TestService_CreateImage(t *testing.T) {
	ctx := context.Background()
	svc := initSvc()
	image, err := svc.CreateImage(ctx, "生成一张小猫的图片", openai.CreateImageSize256x256, openai.CreateImageResponseFormatB64JSON)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(image)
	for _, v := range image {
		t.Log(v.URL)
		t.Log(v.B64JSON)
	}
}

func TestService_Embeddings(t *testing.T) {
	encodedInput := [][]int{
		{37046, 19000, 80172, 73361, 18259, 108, 9554, 75146, 35330, 3922, 17486, 236, 82696, 52254, 198, 80172, 73361, 42052, 33052, 16325, 63344, 21043, 34577, 6701, 116, 22238, 75146, 35330, 9554, 37767, 54493, 47770, 17701, 35287, 45114, 106, 14888, 121, 75146, 35330, 3922, 88126, 74770, 68438, 45114, 106, 14888, 121, 75146, 35330, 15049, 72917, 52254, 3922, 93994, 9554, 42052, 33052, 83800, 165, 118, 119, 163, 225, 99, 26955, 101, 76537, 3443, 12, 24809, 12, 5500, 15, 3922, 99750, 39013, 95, 46456, 69978, 1811, 1432, 47551, 19000, 52225, 52332, 34577, 31958, 17486, 236, 82696, 30867, 198, 88126, 74770, 68438, 45114, 106, 14888, 121, 75146, 35330, 15049, 3922, 73751, 11144, 98711, 11199, 3922, 19000, 60843, 52225, 52332, 9554, 65917, 20135, 100, 19361, 48044, 11144, 52225, 52332, 34577, 31958, 11199, 3922, 88126, 74770, 68438, 33091, 99480, 73751, 64026, 32296, 3922, 15308, 117, 16423, 20379, 27452, 44915, 72917, 69856, 62543, 3922, 51504, 73117, 94668, 34547, 92776, 31540, 1811, 1432, 98711, 75146, 35330, 37767, 54493, 11239, 244, 20834, 35287, 3922, 6271, 222, 82696, 13646, 20022, 247, 74770, 28037, 81201, 198, 16937, 42016, 9554, 75146, 35330, 28425, 236, 18904, 28037, 81201, 9554, 21082, 16937, 42016, 3922, 35417, 53508, 253, 22238, 75146, 35330, 33035, 40053, 86206, 19, 12, 22, 19483, 39209, 87844, 9080, 28037, 81201, 3922, 86127, 53508, 253, 22238, 23897, 40862, 42016, 95337, 23187, 18184, 12870, 228, 1811, 88126, 74770, 73751, 98711, 313, 69978, 6271, 241, 90756, 93233, 28037, 11239, 244, 20834, 9554, 75146, 35330, 91495, 33976, 52030, 99941, 38743, 28037, 81201, 21082, 1811, 1432, 37046, 19361, 48044, 6277, 75146, 35330, 33565, 111, 11239, 244, 19012, 231, 3922, 78388, 10414, 102, 37046, 13079, 98, 37087, 52225, 42506, 14, 22023, 74396, 21043, 6271, 222, 82696, 37689, 91763, 198, 88126, 53901, 3922, 75146, 35330, 28425, 236, 18904, 25580, 15308, 117, 16423, 49928, 36651, 9554, 31634, 32018, 88126, 86206, 47406, 83601, 226, 18259, 251, 22238, 31634, 72238, 9554, 28469, 3922, 94668, 21043, 92019, 48044, 59614, 69856, 48982, 3922, 88126, 69856, 62543, 48044, 74770, 51109, 28037, 9554, 94668, 47577, 3922, 69905, 73597, 94, 50266, 111, 4239, 31, 4239, 95429, 92776, 31540, 10110, 30624, 90651, 33721, 62752, 916, 8},
		{28194, 117, 22649, 13647, 95, 14, 75146, 35330, 44510, 228, 14, 53901, 75486, 6701, 240, 14, 80172, 73361, 5726, 14, 8676, 250, 17792, 14, 78519, 80699, 14, 13647, 95, 15973, 234, 14, 41127, 40862, 92553, 17982, 55030, 34547, 17486, 236, 82696, 28037, 35287, 45114, 106, 14888, 121, 75146, 35330, 3922, 37046, 31634, 72917, 50338, 12774, 239, 198, 61075, 61826, 28194, 117, 22649, 13647, 95, 9554, 75146, 35330, 22238, 83800, 15308, 117, 16423, 49928, 36651, 31634, 32018, 37767, 54493, 64026, 17701, 28037, 45114, 106, 14888, 121, 75146, 35330, 77413, 20379, 3922, 63344, 88126, 30832, 25580, 81543, 69978, 6271, 241, 9554, 52225, 52332, 74770, 7518, 111, 15225, 81201, 17982, 26130, 92553, 3922, 88126, 31540, 19000, 45114, 106, 14888, 121, 75146, 35330, 15049, 73751, 11144, 98711, 11199, 313, 11144, 45018, 11199, 313, 11144, 34226, 43240, 11199, 313, 11144, 26130, 92553, 81201, 17982, 11199, 37507, 72917, 41190, 3922, 63344, 86206, 26130, 92553, 64531, 19483, 75146, 35330, 81201, 17982, 3922, 88126, 74770, 68438, 3443, 12, 21138, 24, 12, 1049, 58318, 98739, 73164, 1811, 1432, 37046, 19000, 8676, 250, 17792, 13647, 115, 7518, 111, 15225, 9554, 13647, 115, 69253, 48974, 3922, 37046, 31634, 52030, 17297, 12774, 228, 56602, 198, 88126, 53901, 3922, 66776, 40053, 24326, 109, 15722, 231, 98739, 44388, 70349, 21043, 45114, 106, 14888, 121, 75146, 35330, 3922, 13647, 253, 70616, 34577, 6701, 116, 22238, 75146, 35330, 9554, 80172, 52225, 48974, 3922, 63344, 21043, 20022, 253, 69253, 88126, 26955, 101, 76537, 3443, 12, 22588, 12, 14261, 23, 163, 225, 255, 44368, 161, 240, 101, 39365, 3922, 39013, 95, 39013, 95, 22649, 50338, 5940, 1432, 18184, 6271, 222, 82696, 45114, 106, 14888, 121, 75146, 35330, 70349, 28190, 9554, 72368, 43292, 25333, 18259, 108, 11239, 244, 198, 45114, 106, 14888, 121, 75146, 35330, 30832, 25580, 74770, 72917, 28469, 52254, 5486, 69978, 6271, 241, 52254, 5486, 28425, 236, 18904, 11239, 244, 20834, 3922, 88399, 224, 13646, 98806, 54253, 7518, 111, 78272, 3922, 24326, 109, 15722, 231, 90112, 88126, 13821, 99, 37507, 16937, 20135, 123, 28308, 99, 1811, 77913, 33014, 7518, 111, 78272, 21082, 98739, 86206, 50667, 75293, 32335, 12774, 230, 9554, 91951, 60251, 3922, 75863, 38093, 68438, 163, 253, 255, 22023, 5486, 71890, 32943, 65305, 33035, 53283, 88126, 35287, 50338, 3922, 39013, 95, 39013, 95, 1811},
	}
	var allInput []int
	for _, val := range encodedInput {
		for _, num := range val {
			allInput = append(allInput, num)
		}
	}
	tk, err := tiktoken.EncodingForModel("text-embedding-ada-002")
	ctx := context.Background()
	svc := initSvc()
	embeddings, err := svc.Embeddings(ctx, "text-embedding-ada-002", tk.Decode(allInput))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(embeddings.Usage)
}

func TestService_UploadFile(t *testing.T) {
	ctx := context.Background()
	svc := initSvc()
	res, err := svc.UploadFile(ctx, "gpt-3.5-turbo", "ft.jsonl", "github.com/IceBearAI/aigc/src/services/fastchat/ft.jsonl", "fine-tune")
	if err != nil {
		t.Error(err)
		return
	}
	b, _ := json.Marshal(res)
	t.Log(string(b))
	// {"bytes":157495,"created_at":1693386923,"id":"file-CYLNxI7ppCkKDiGxCpL5X5zq","filename":"github.com/IceBearAI/aigc/src/services/fastchat/ft.jsonl","object":"file","owner":"","purpose":"fine-tune"}
}

func TestService_CreateFineTuningJob(t *testing.T) {
	ctx := context.Background()
	svc := initSvc()
	res, err := svc.CreateFineTuningJob(ctx, openai.FineTuningJobRequest{
		TrainingFile:   "file-CYLNxI7ppCkKDiGxCpL5X5zq",
		ValidationFile: "",
		Model:          "gpt-3.5-turbo",
		Hyperparameters: &openai.Hyperparameters{
			Epochs: 3,
		},
		Suffix: "",
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(res)
	// {"id":"ftjob-hbF6RLCsVlS48UrtQt9qTlSd","object":"fine_tuning.job","created_at":1693387245,"finished_at":0,"model":"gpt-3.5-turbo-0613","organization_id":"org-CKKN6woaf472fyl3c3zBV6JA","status":"created","hyperparameters":{"n_epochs":3},"training_file":"file-CYLNxI7ppCkKDiGxCpL5X5zq","result_files":[],"trained_tokens":0}
}

func TestService_RetrieveFineTuningJob(t *testing.T) {
	ctx := context.Background()
	svc := initSvc()
	res, err := svc.RetrieveFineTuningJob(ctx, "ftjob-hbF6RLCsVlS48UrtQt9qTlSd")
	if err != nil {
		t.Error(err)
		return
	}
	b, _ := json.Marshal(res)
	t.Log(string(b))
	// {"object":"fine_tuning.job","id":"ftjob-hbF6RLCsVlS48UrtQt9qTlSd","model":"gpt-3.5-turbo-0613","created_at":1693387245,"finished_at":null,"fine_tuned_model":null,"organization_id":"org-CKKN6woaf472fyl3c3zBV6JA","result_files":[],"status":"running","validation_file":null,"training_file":"file-CYLNxI7ppCkKDiGxCpL5X5zq","hyperparameters":{"n_epochs":3},"trained_tokens":null}
}

func TestService_ListFineTune(t *testing.T) {
	ctx := context.Background()
	svc := initSvc()
	res, err := svc.ListFineTune(ctx, "gpt-3.5-turbo")
	if err != nil {
		t.Error(err)
		return
	}
	b, _ := json.Marshal(res)
	t.Log(string(b))
}

func TestService_CancelFineTuningJob(t *testing.T) {
	ctx := context.Background()
	svc := initSvc()
	err := svc.CancelFineTuningJob(ctx, "gpt-3.5-turbo", "ftjob-Va7GX0tDc3ABOQlVbE76lHBu")
	if err != nil {
		t.Error(err)
		return
	}
}
