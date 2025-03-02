package main

import (
	"fmt"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

func synthesizeStartedHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	fmt.Println("Synthesis started.")
}

func synthesizingHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	fmt.Printf("Synthesizing, audio chunk size %d.\n", len(event.Result.AudioData))
}

func synthesizedHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	fmt.Printf("Synthesized, audio length %d.\n", len(event.Result.AudioData))
}

func cancelledHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	fmt.Println("Received a cancellation.")
}

func main() {
	// 設定 Azure Speech API 金鑰與區域
	speechKey := "your_speech_key" // 請填入你的 Azure 訂閱金鑰
	speechRegion := "your_region"  // 例如 "eastus"

	// **不設定音訊輸出，避免 ALSA 錯誤**
	var audioConfig *audio.AudioConfig = nil

	// 建立語音設定
	speechConfig, err := speech.NewSpeechConfigFromSubscription(speechKey, speechRegion)
	if err != nil {
		fmt.Println("Speech config error:", err)
		return
	}
	defer speechConfig.Close()

	// 設定語音合成的語音
	speechConfig.SetSpeechSynthesisVoiceName("en-US-AvaMultilingualNeural")

	// 建立語音合成器（不播放音訊）
	speechSynthesizer, err := speech.NewSpeechSynthesizerFromConfig(speechConfig, audioConfig)
	if err != nil {
		fmt.Println("Speech synthesizer error:", err)
		return
	}
	defer speechSynthesizer.Close()

	// 設定事件監聽
	speechSynthesizer.SynthesisStarted(synthesizeStartedHandler)
	speechSynthesizer.Synthesizing(synthesizingHandler)
	speechSynthesizer.SynthesisCompleted(synthesizedHandler)
	speechSynthesizer.SynthesisCanceled(cancelledHandler)

	// 需要發音的文字
	text := "Hello, welcome to Azure Speech SDK!"

	// 執行語音合成
	task := speechSynthesizer.SpeakTextAsync(text)
	var outcome speech.SpeechSynthesisOutcome

	// 設定超時 30 秒
	select {
	case outcome = <-task:
	case <-time.After(30 * time.Second):
		fmt.Println("Speech synthesis timed out")
		return
	}
	defer outcome.Close()

	// 檢查錯誤
	if outcome.Error != nil {
		fmt.Println("Speech synthesis error:", outcome.Error)
		return
	}

	// 成功訊息
	if outcome.Result.Reason == common.SynthesizingAudioCompleted {
		fmt.Printf("Speech synthesized successfully for text [%s].\n", text)
	} else {
		cancellation, _ := speech.NewCancellationDetailsFromSpeechSynthesisResult(outcome.Result)
		fmt.Printf("CANCELED: Reason=%d.\n", cancellation.Reason)

		if cancellation.Reason == common.Error {
			fmt.Printf("CANCELED: ErrorCode=%d\nCANCELED: ErrorDetails=[%s]\nCANCELED: Did you set the speech resource key and region values?\n",
				cancellation.ErrorCode,
				cancellation.ErrorDetails)
		}
	}
}
