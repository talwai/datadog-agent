// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

// +build python

package python

import (
	"unsafe"

	"github.com/DataDog/datadog-agent/pkg/tagger"
	"github.com/DataDog/datadog-agent/pkg/tagger/collectors"
)

/*
#include <datadog_agent_six.h>
#cgo !windows LDFLAGS: -ldatadog-agent-six -ldl
#cgo windows LDFLAGS: -ldatadog-agent-six -lstdc++ -static

*/
import "C"

// Tags bridges towards tagger.Tag to retrieve container tags
//export Tags
func Tags(id *C.char, cardinality C.int) **C.char {
	goID := C.GoString(id)
	var tags []string

	tags, _ = tagger.Tag(goID, collectors.TagCardinality(cardinality))

	length := len(tags)
	if length == 0 {
		return nil
	}

	cTags := C.malloc(C.size_t(length+1) * C.size_t(unsafe.Sizeof(uintptr(0))))

	// convert the C array to a Go Array so we can index it
	indexTag := (*[1<<29 - 1]*C.char)(cTags)[: length+1 : length+1]
	indexTag[length] = nil
	for idx, tag := range tags {
		indexTag[idx] = C.CString(tag)
	}

	return (**C.char)(cTags)
}
