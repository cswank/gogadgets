package io

type OutputDevice interface {
	On() error
	Off() error
}


