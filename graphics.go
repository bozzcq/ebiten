// Copyright 2014 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ebiten

import (
	"github.com/hajimehoshi/ebiten/internal/graphics"
	"github.com/hajimehoshi/ebiten/internal/opengl"
)

// Filter represents the type of texture filter to be used when an image is maginified or minified.
type Filter int

const (
	// FilterDefault represents the defualt filter.
	FilterDefault Filter = Filter(graphics.FilterDefault)

	// FilterNearest represents nearest (crisp-edged) filter
	FilterNearest Filter = Filter(graphics.FilterNearest)

	// FilterLinear represents linear filter
	FilterLinear Filter = Filter(graphics.FilterLinear)
)

// CompositeMode represents Porter-Duff composition mode.
type CompositeMode int

// This name convention follows CSS compositing: https://drafts.fxtf.org/compositing-2/.
//
// In the comments,
// c_src, c_dst and c_out represent alpha-premultiplied RGB values of source, destination and output respectively. α_src and α_dst represent alpha values of source and destination respectively.
const (
	// Regular alpha blending
	// c_out = c_src + c_dst × (1 - α_src)
	CompositeModeSourceOver CompositeMode = CompositeMode(opengl.CompositeModeSourceOver)

	// c_out = 0
	CompositeModeClear CompositeMode = CompositeMode(opengl.CompositeModeClear)

	// c_out = c_src
	CompositeModeCopy CompositeMode = CompositeMode(opengl.CompositeModeCopy)

	// c_out = c_dst
	CompositeModeDestination CompositeMode = CompositeMode(opengl.CompositeModeDestination)

	// c_out = c_src × (1 - α_dst) + c_dst
	CompositeModeDestinationOver CompositeMode = CompositeMode(opengl.CompositeModeDestinationOver)

	// c_out = c_src × α_dst
	CompositeModeSourceIn CompositeMode = CompositeMode(opengl.CompositeModeSourceIn)

	// c_out = c_dst × α_src
	CompositeModeDestinationIn CompositeMode = CompositeMode(opengl.CompositeModeDestinationIn)

	// c_out = c_src × (1 - α_dst)
	CompositeModeSourceOut CompositeMode = CompositeMode(opengl.CompositeModeSourceOut)

	// c_out = c_dst × (1 - α_src)
	CompositeModeDestinationOut CompositeMode = CompositeMode(opengl.CompositeModeDestinationOut)

	// c_out = c_src × α_dst + c_dst × (1 - α_src)
	CompositeModeSourceAtop CompositeMode = CompositeMode(opengl.CompositeModeSourceAtop)

	// c_out = c_src × (1 - α_dst) + c_dst × α_src
	CompositeModeDestinationAtop CompositeMode = CompositeMode(opengl.CompositeModeDestinationAtop)

	// c_out = c_src × (1 - α_dst) + c_dst × (1 - α_src)
	CompositeModeXor CompositeMode = CompositeMode(opengl.CompositeModeXor)

	// Sum of source and destination (a.k.a. 'plus' or 'additive')
	// c_out = c_src + c_dst
	CompositeModeLighter CompositeMode = CompositeMode(opengl.CompositeModeLighter)
)
