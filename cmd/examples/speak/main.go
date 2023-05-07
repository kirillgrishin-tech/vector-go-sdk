package main

import (
	"context"
	"flag"
	"fmt"

	sdk_wrapper "github.com/kirillgrishin-tech/vector-go-sdk/pkg/sdk-wrapper"
)

func main() {
	var serial = flag.String("serial", "", "Vector's Serial Number")
	flag.Parse()

	sdk_wrapper.InitSDK(*serial)

	ctx := context.Background()
	start := make(chan bool)
	stop := make(chan bool)

	go func() {
		_ = sdk_wrapper.Robot.BehaviorControl(ctx, start, stop)
	}()

	for {
		select {
		case <-start:
			/*sdk_wrapper.SayText("I wouldn't have let anything come between")
			You wouldn't have done that if there had not been something
			"he's afraid of being seen, being found"
			"her maid was standng by the garden gate, looking for her"
			"she had done all that was possible"
			 "it was the seal upon the bond"
			"the odds between her and her adversary were even"
			"she would have to break her word to Milly"
			"she had a light burning in her room till morning, for she was afraid of sleep"
			"her gift, her secret, was powerless now against the pursuer"
			"a terrified bird flew out of the hedge, no further than a flight in front of her"
			"all this she perceived in a flash, when she had turned the corner"
			"as she turned the corner of the wood, she was brought suddenly in sight of the valley"
			"now her fear, which had become almost hatred, was transferred to his person"
			"what she saw was the empty shell of him"
			"she went on and came to the gate of the wood"
			"she paused on the bridge, and looked down the valley"
			"it was perfect, following a perfect day"
			"she waited for her hour between sunset and twilight"
			 "she told herself that, after all, her fear had done no harm"
			"she was bound to accept his statement"
			"she doesn't care a rap about me"
			"they had sat down on the couch in the corner so that they faced each other"
			"she begged him to write and tell her that he was well"
			"she refused to hold him even by a thread"
			*/
			/*
				sdk_wrapper.SetLocale("en-US")
				sdk_wrapper.SetTTSEngine(sdk_wrapper.TTS_ENGINE_NATIVE)
				sdk_wrapper.SayText("Hello world. I really want to explore!")
				sdk_wrapper.SetTTSEngine(sdk_wrapper.TTS_ENGINE_ESPEAK)
				sdk_wrapper.SayText("Hello world. I really want to explore!")
				sdk_wrapper.SetTTSEngine(sdk_wrapper.TTS_ENGINE_HTGO)
				sdk_wrapper.SayText("Hello world. I really want to explore!")
				sdk_wrapper.SetTTSEngine(sdk_wrapper.TTS_ENGINE_VOICESERVER)
				sdk_wrapper.SetTTSVoice(sdk_wrapper.TTS_ENGINE_VOICESERVER_VOICE_ENGLISH_WEDNESDAY)
				sdk_wrapper.SayText("Hello world. I really want to explore!")
				sdk_wrapper.SetTTSVoice(sdk_wrapper.TTS_ENGINE_VOICESERVER_VOICE_ENGLISH_VOLDEMORT)
				sdk_wrapper.SayText("Hello world. I really want to explore!")
				sdk_wrapper.SetTTSVoice(sdk_wrapper.TTS_ENGINE_VOICESERVER_VOICE_ENGLISH_BATTLEDROID)
				sdk_wrapper.SayText("Hello world. I really want to explore!")

				sdk_wrapper.SetLocale("it-IT")
				sdk_wrapper.SetTTSEngine(sdk_wrapper.TTS_ENGINE_NATIVE)
				sdk_wrapper.SayText("Ciao mondo! Sono pronto per esplorare")
				sdk_wrapper.SetTTSEngine(sdk_wrapper.TTS_ENGINE_ESPEAK)
				sdk_wrapper.SayText("Ciao mondo! Sono pronto per esplorare")
				sdk_wrapper.SetTTSEngine(sdk_wrapper.TTS_ENGINE_HTGO)
				sdk_wrapper.SayText("Ciao mondo! Sono pronto per esplorare")
				sdk_wrapper.SetTTSEngine(sdk_wrapper.TTS_ENGINE_VOICESERVER)
				sdk_wrapper.SetTTSVoice(sdk_wrapper.TTS_ENGINE_VOICESERVER_VOICE_ITALIAN_WANNA_MARCHI)
				sdk_wrapper.SayText("Ciao mondo! Sono pronto per esplorare")
			*/

			sdk_wrapper.SetLocale("de-DE")
			fmt.Println(sdk_wrapper.GetLocale())
			sdk_wrapper.SetTTSEngine(sdk_wrapper.TTS_ENGINE_NATIVE)
			sdk_wrapper.SayText("Hallo Welt! Ich bin bereit zu erkunden")
			/*
				sdk_wrapper.SetTTSEngine(sdk_wrapper.TTS_ENGINE_ESPEAK)
				sdk_wrapper.SayText("Hallo Welt! Ich bin bereit zu erkunden")
				sdk_wrapper.SetTTSEngine(sdk_wrapper.TTS_ENGINE_HTGO)
				sdk_wrapper.SayText("Hallo Welt! Ich bin bereit zu erkunden")
				sdk_wrapper.SetTTSEngine(sdk_wrapper.TTS_ENGINE_VOICESERVER)
				sdk_wrapper.SetTTSVoice(sdk_wrapper.TTS_ENGINE_VOICESERVER_VOICE_GERMAN_WOLFF_FUSS)
				sdk_wrapper.SayText("Hallo Welt! Ich bin bereit zu erkunden")

				sdk_wrapper.SetLocale("fr-FR")
				sdk_wrapper.SetTTSEngine(sdk_wrapper.TTS_ENGINE_NATIVE)
				sdk_wrapper.SayText("Bonjour le monde! Je suis prêt à explorer")
				sdk_wrapper.SetTTSEngine(sdk_wrapper.TTS_ENGINE_ESPEAK)
				sdk_wrapper.SayText("Bonjour le monde! Je suis prêt à explorer")
				sdk_wrapper.SetTTSEngine(sdk_wrapper.TTS_ENGINE_HTGO)
				sdk_wrapper.SayText("Bonjour le monde! Je suis prêt à explorer")
				sdk_wrapper.SetTTSEngine(sdk_wrapper.TTS_ENGINE_VOICESERVER)
				sdk_wrapper.SetTTSVoice(sdk_wrapper.TTS_ENGINE_VOICESERVER_VOICE_FRENCH_BUGS_BUNNY)
				sdk_wrapper.SayText("Bonjour le monde! Je suis prêt à explorer")

				sdk_wrapper.SetLocale("es-ES")
				sdk_wrapper.SetTTSEngine(sdk_wrapper.TTS_ENGINE_NATIVE)
				sdk_wrapper.SayText("¡Hola Mundo! Estoy listo para explorar")
				sdk_wrapper.SetTTSEngine(sdk_wrapper.TTS_ENGINE_ESPEAK)
				sdk_wrapper.SayText("¡Hola Mundo! Estoy listo para explorar")
				sdk_wrapper.SetTTSEngine(sdk_wrapper.TTS_ENGINE_HTGO)
				sdk_wrapper.SayText("¡Hola Mundo! Estoy listo para explorar")
				sdk_wrapper.SetTTSEngine(sdk_wrapper.TTS_ENGINE_VOICESERVER)
				sdk_wrapper.SetTTSVoice(sdk_wrapper.TTS_ENGINE_VOICESERVER_VOICE_SPANISH_BURRO_DE_SHREK)
				sdk_wrapper.SayText("¡Hola Mundo! Estoy listo para explorar")
			*/
			stop <- true
			return
		}
	}
}
