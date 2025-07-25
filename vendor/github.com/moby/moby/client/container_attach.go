package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/moby/moby/api/types"
	"github.com/moby/moby/api/types/container"
)

// ContainerAttach attaches a connection to a container in the server.
// It returns a types.HijackedConnection with the hijacked connection
// and the a reader to get output. It's up to the called to close
// the hijacked connection by calling types.HijackedResponse.Close.
//
// The stream format on the response will be in one of two formats:
//
// If the container is using a TTY, there is only a single stream (stdout), and
// data is copied directly from the container output stream, no extra
// multiplexing or headers.
//
// If the container is *not* using a TTY, streams for stdout and stderr are
// multiplexed.
// The format of the multiplexed stream is as follows:
//
//	[8]byte{STREAM_TYPE, 0, 0, 0, SIZE1, SIZE2, SIZE3, SIZE4}[]byte{OUTPUT}
//
// STREAM_TYPE can be 1 for stdout and 2 for stderr
//
// SIZE1, SIZE2, SIZE3, and SIZE4 are four bytes of uint32 encoded as big endian.
// This is the size of OUTPUT.
//
// You can use github.com/moby/moby/api/stdcopy.StdCopy to demultiplex this
// stream.
func (cli *Client) ContainerAttach(ctx context.Context, containerID string, options container.AttachOptions) (types.HijackedResponse, error) {
	containerID, err := trimID("container", containerID)
	if err != nil {
		return types.HijackedResponse{}, err
	}

	query := url.Values{}
	if options.Stream {
		query.Set("stream", "1")
	}
	if options.Stdin {
		query.Set("stdin", "1")
	}
	if options.Stdout {
		query.Set("stdout", "1")
	}
	if options.Stderr {
		query.Set("stderr", "1")
	}
	if options.DetachKeys != "" {
		query.Set("detachKeys", options.DetachKeys)
	}
	if options.Logs {
		query.Set("logs", "1")
	}

	return cli.postHijacked(ctx, "/containers/"+containerID+"/attach", query, nil, http.Header{
		"Content-Type": {"text/plain"},
	})
}
