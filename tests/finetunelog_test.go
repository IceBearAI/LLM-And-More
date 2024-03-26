package tests

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"
)

type logEntry struct {
	Timestamp    time.Time `json:"timestamp"`
	Loss         float64   `json:"loss"`
	LearningRate float64   `json:"learning_rate"`
	Epoch        float64   `json:"epoch"`
	GradNorm     float64   `json:"grad_norm"`
}

func TestFinetuneLog(t *testing.T) {

	jobLog := `/usr/local/lib/python3.10/dist-packages/torch/autograd/__init__.py:266: UserWarning: c10d::broadcast_: an autograd kernel was not registered to the Autograd key(s) but we are trying to backprop through it. This may lead to silently incorrect behavior. This behavior is deprecated and will be removed in a future version of PyTorch. If your operator is differentiable, please ensure you have registered an autograd kernel to the correct Autograd key (e.g. DispatchKey::Autograd, DispatchKey::CompositeImplicitAutograd). If your operator is not differentiable, or to squash this warning and use the previous behavior, please register torch::CppFunction::makeFallthrough() to DispatchKey::Autograd. (Triggered internally at ../torch/csrc/autograd/autograd_not_implemented_fallback.cpp:63.)
a  Variable._execution_engine.run_backward(  # Calls into the C++ engine to run the backward pass
/usr/local/lib/python3.10/dist-packages/torch/autograd/__init__.py:266: UserWarning: c10d::broadcast_: an autograd kernel was not registered to the Autograd key(s) but we are trying to backprop through it. This may lead to silently incorrect behavior. This behavior is deprecated and will be removed in a future version of PyTorch. If your operator is differentiable, please ensure you have registered an autograd kernel to the correct Autograd key (e.g. DispatchKey::Autograd, DispatchKey::CompositeImplicitAutograd). If your operator is not differentiable, or to squash this warning and use the previous behavior, please register torch::CppFunction::makeFallthrough() to DispatchKey::Autograd. (Triggered internally at ../torch/csrc/autograd/autograd_not_implemented_fallback.cpp:63.)
a  Variable._execution_engine.run_backward(  # Calls into the C++ engine to run the backward pass
V{'loss': 3.824, 'grad_norm': 1.7319538297027568, 'learning_rate': 0.0, 'epoch': 0.33}
{'loss': 3.8254, 'grad_norm': 1.770027949872221, 'learning_rate': 2e-05, 'epoch': 0.67}
X{'loss': 3.8314, 'grad_norm': 1.8041480112616333, 'learning_rate': 2e-05, 'epoch': 1.0}
{'loss': 3.8095, 'grad_norm': 1.7225661481272851, 'learning_rate': 2e-05, 'epoch': 1.33}
{'loss': 2.4368, 'grad_norm': 0.7307336908165509, 'learning_rate': 2e-05, 'epoch': 19.33}
Z{'loss': 2.4019, 'grad_norm': 0.7154495507267877, 'learning_rate': 2e-05, 'epoch': 19.67}
 88%|████████▊ 2024-03-06T15:48:30.364072185Z {'loss': 0.0532, 'learning_rate': 2e-05, 'epoch': 4.39}
2024-03-06T15:48:55.232199457Z {'loss': 0.063, 'learning_rate': 2e-05, 'epoch': 4.41}
2024-03-06T15:49:20.085236223Z {'loss': 0.0683, 'learning_rate': 2e-05, 'epoch': 4.42}
2024-03-06T15:49:44.949303974Z {'loss': 0.0627, 'learning_rate': 2e-05, 'epoch': 4.44}
2024-03-06T15:50:09.814279635Z {'loss': 0.0585, 'learning_rate': 2e-05, 'epoch': 4.45}
2024-03-06T15:50:59.519041650Z {'loss': 0.0789, 'learning_rate': 2e-05, 'epoch': 4.48}
[1612] 2024-03-06T15:51:49.223750867Z {'loss': 0.076, 'learning_rate': 2e-05, 'epoch': 4.52}`

	lineArr := strings.Split(jobLog, "\n")
	// 正则表达式匹配花括号内的内容
	jsonPattern := regexp.MustCompile(`\{[^}]*\}`)

	var logEntryList []logEntry

	for _, log := range lineArr {
		log = strings.TrimSpace(log)
		matches := jsonPattern.FindAllString(log, -1)
		for _, match := range matches {
			if len(match) > 0 {
				// 将单引号替换为双引号以符合JSON格式
				jsonStr := strings.Replace(match, "'", "\"", -1)         // 将单引号替换为双引号
				jsonStr = strings.Replace(jsonStr, "False", "false", -1) // 将 False 替换为 false
				jsonStr = strings.Replace(jsonStr, "True", "true", -1)   // 将 True 替换为 true

				var entry logEntry
				err := json.Unmarshal([]byte(jsonStr), &entry)
				if err != nil {
					fmt.Printf("unmarshal json failed: %s\n", err.Error())
					continue
				}

				logEntryList = append(logEntryList, entry)
			}
		}
	}

	for _, v := range logEntryList {
		b, _ := json.Marshal(v)
		t.Log(string(b))
	}
}
