package fastchat

type Txt2Img struct {
	EnableHr          bool     `json:"enable_hr"`          //是否启用高分辨率
	DenoisingStrength int      `json:"denoising_strength"` //降噪强度
	FirstphaseWidth   int      `json:"firstphase_width"`
	FirstphaseHeight  int      `json:"firstphase_height"`
	HrScale           int      `json:"hr_scale"`             //高分辨率缩放比例
	HrUpscaler        string   `json:"hr_upscaler"`          //高分辨率缩放的方法或算法
	HrSecondPassSteps int      `json:"hr_second_pass_steps"` //高分辨率生成过程的迭代步数
	HrResizeX         int      `json:"hr_resize_x"`          //高分辨 宽
	HrResizeY         int      `json:"hr_resize_y"`          //高分辨 高
	HrSamplerName     string   `json:"hr_sampler_name"`      //高分辨 采样器
	HrPrompt          string   `json:"hr_prompt"`            //高分辨 提示词 正向
	HrNegativePrompt  string   `json:"hr_negative_prompt"`   //高分辨 提示词 负向
	Prompt            string   `json:"prompt"`               //提示词 正向
	Styles            []string `json:"styles"`               //风格
	Seed              int      `json:"seed"`                 //种子
	Subseed           int      `json:"subseed"`              //种子 进一步调整算法中的随机性
	SubseedStrength   int      `json:"subseed_strength"`     //种子强度
	SeedResizeFromH   int      `json:"seed_resize_from_h"`   //从哪个高度开始调整种子
	SeedResizeFromW   int      `json:"seed_resize_from_w"`   //从哪个宽度开始调整种子
	SamplerName       string   `json:"sampler_name"`         //采样器
	BatchSize         int      `json:"batch_size"`           //批次
	NIter             int      `json:"n_iter"`               //迭代次数
	Steps             int      `json:"steps"`                //步数
	CfgScale          int      `json:"cfg_scale"`            //缩放因子
	Width             int      `json:"width"`
	Height            int      `json:"height"`
	RestoreFaces      bool     `json:"restore_faces"`
	Tiling            bool     `json:"tiling"`
	DoNotSaveSamples  bool     `json:"do_not_save_samples"`
	DoNotSaveGrid     bool     `json:"do_not_save_grid"`
	NegativePrompt    string   `json:"negative_prompt"` //提示词 负向
	Eta               int      `json:"eta"`
	SMinUncond        int      `json:"s_min_uncond"`
	SChurn            int      `json:"s_churn"`
	STmax             int      `json:"s_tmax"`
	STmin             int      `json:"s_tmin"`
	SNoise            int      `json:"s_noise"`
	OverrideSettings  struct {
	} `json:"override_settings"`
	OverrideSettingsRestoreAfterwards bool          `json:"override_settings_restore_afterwards"`
	ScriptArgs                        []interface{} `json:"script_args"`
	SamplerIndex                      string        `json:"sampler_index"`
	ScriptName                        string        `json:"script_name"`
	SendImages                        bool          `json:"send_images"`
	SaveImages                        bool          `json:"save_images"`
	AlwaysonScripts                   struct {
	} `json:"alwayson_scripts"`
}

// sd
type SdSuccessRep struct {
	Images     []string `json:"images"`
	Parameters struct {
		EnableHr                          bool          `json:"enable_hr"`
		DenoisingStrength                 int           `json:"denoising_strength"`
		FirstphaseWidth                   int           `json:"firstphase_width"`
		FirstphaseHeight                  int           `json:"firstphase_height"`
		HrScale                           float64       `json:"hr_scale"`
		HrUpscaler                        interface{}   `json:"hr_upscaler"`
		HrSecondPassSteps                 int           `json:"hr_second_pass_steps"`
		HrResizeX                         int           `json:"hr_resize_x"`
		HrResizeY                         int           `json:"hr_resize_y"`
		HrSamplerName                     interface{}   `json:"hr_sampler_name"`
		HrPrompt                          string        `json:"hr_prompt"`
		HrNegativePrompt                  string        `json:"hr_negative_prompt"`
		Prompt                            string        `json:"prompt"`
		Styles                            interface{}   `json:"styles"`
		Seed                              int           `json:"seed"`
		Subseed                           int           `json:"subseed"`
		SubseedStrength                   int           `json:"subseed_strength"`
		SeedResizeFromH                   int           `json:"seed_resize_from_h"`
		SeedResizeFromW                   int           `json:"seed_resize_from_w"`
		SamplerName                       interface{}   `json:"sampler_name"`
		BatchSize                         int           `json:"batch_size"`
		NIter                             int           `json:"n_iter"`
		Steps                             int           `json:"steps"`
		CfgScale                          float64       `json:"cfg_scale"`
		Width                             int           `json:"width"`
		Height                            int           `json:"height"`
		RestoreFaces                      bool          `json:"restore_faces"`
		Tiling                            bool          `json:"tiling"`
		DoNotSaveSamples                  bool          `json:"do_not_save_samples"`
		DoNotSaveGrid                     bool          `json:"do_not_save_grid"`
		NegativePrompt                    interface{}   `json:"negative_prompt"`
		Eta                               interface{}   `json:"eta"`
		SMinUncond                        float64       `json:"s_min_uncond"`
		SChurn                            float64       `json:"s_churn"`
		STmax                             interface{}   `json:"s_tmax"`
		STmin                             float64       `json:"s_tmin"`
		SNoise                            float64       `json:"s_noise"`
		OverrideSettings                  interface{}   `json:"override_settings"`
		OverrideSettingsRestoreAfterwards bool          `json:"override_settings_restore_afterwards"`
		ScriptArgs                        []interface{} `json:"script_args"`
		SamplerIndex                      string        `json:"sampler_index"`
		ScriptName                        interface{}   `json:"script_name"`
		SendImages                        bool          `json:"send_images"`
		SaveImages                        bool          `json:"save_images"`
		AlwaysonScripts                   struct {
		} `json:"alwayson_scripts"`
	} `json:"parameters"`
	Info string `json:"info"`
}

// ImageProgress sd 生成过程 progress Result
type ImageProgress struct {
	Progress    float64 `json:"progress"` //进度百分比 完成后为 "0.0"
	EtaRelative float64 `json:"eta_relative"`
	State       struct {
		Skipped       bool   `json:"skipped"`     //跳过
		Interrupted   bool   `json:"interrupted"` //中止
		Job           string `json:"job"`
		JobCount      int    `json:"job_count"`
		JobTimestamp  string `json:"job_timestamp"`  //当前时间 格式 "20230628124125"
		JobNo         int    `json:"job_no"`         //过程中=0，结束后=1
		SamplingStep  int    `json:"sampling_step"`  //当前步数
		SamplingSteps int    `json:"sampling_steps"` //总步数
	} `json:"state"`
	CurrentImage string      `json:"current_image"` //图片信息
	Textinfo     interface{} `json:"textinfo"`
}

// Txt2ImgRequest 文生图 Request
type Txt2ImgRequest struct {
	Steps          int    `json:"steps"`
	Prompt         string `json:"prompt"`          //正向 提示词ImageProgress
	NegativePrompt string `json:"negative_prompt"` //负向 提示词
	SamplerIndex   string `json:"sampler_index"`   //采样方法
}

// Txt2ImgResult 文生图 Result
type Txt2ImgResult struct {
	Finish        bool          `json:"finish"`
	Error         error         `json:"error"`
	ImageProgress ImageProgress `json:"imageProgress"`
}

// CreateSdImage 生成图片 Result
type CreateSdImage struct {
	Error        error        `json:"error"`
	SdSuccessRep SdSuccessRep `json:"sdSuccessRep"`
}

type EmbeddingResponse struct {
	Object string          `json:"object"`
	Data   []EmbeddingData `json:"data"`
	Model  string          `json:"model"`
	Usage  EmbeddingUsage  `json:"usage"`
}

type EmbeddingData struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

type EmbeddingUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}
