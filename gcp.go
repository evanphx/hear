package hear

import (
	"context"

	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	"google.golang.org/grpc"

	gcontext "golang.org/x/net/context"
	speech "google.golang.org/genproto/googleapis/cloud/speech/v1beta1"
)

type GCPSpeechConv struct {
	ctx  gcontext.Context
	conn *grpc.ClientConn

	client speech.SpeechClient
}

func NewGCPSpeechConv(accountFile string) (*GCPSpeechConv, error) {
	ctx := context.Background()
	conn, err := transport.DialGRPC(ctx,
		option.WithEndpoint("speech.googleapis.com:443"),
		option.WithScopes("https://www.googleapis.com/auth/cloud-platform"),
		option.WithServiceAccountFile(accountFile),
	)

	if err != nil {
		return nil, err
	}

	client := speech.NewSpeechClient(conn)

	return &GCPSpeechConv{ctx, conn, client}, nil
}

func (gcp *GCPSpeechConv) Convert(data []byte) (string, error) {
	resp, err := gcp.recognize(data)
	if err != nil {
		return "", err
	}

	var best *speech.SpeechRecognitionAlternative

	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			if best == nil || alt.Confidence > best.Confidence {
				best = alt
			}
		}
	}

	return best.Transcript, nil
}

func (gcp *GCPSpeechConv) recognize(data []byte) (*speech.SyncRecognizeResponse, error) {
	return gcp.client.SyncRecognize(gcp.ctx, &speech.SyncRecognizeRequest{
		Config: &speech.RecognitionConfig{
			Encoding:   speech.RecognitionConfig_LINEAR16,
			SampleRate: 16000,
		},
		Audio: &speech.RecognitionAudio{
			AudioSource: &speech.RecognitionAudio_Content{Content: data},
		},
	})
}
