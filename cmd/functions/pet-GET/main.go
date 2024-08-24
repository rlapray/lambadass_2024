//go:build !exclude

package main

import (
	. "github.com/lambadass-2024/backend/cmd/functions/pet-GET/handler"
)

func main() {
	Lambda.
		Use(&Logger).
		Use(&Lambda).
		Use(&SQL).
		Use(&PetRepository).
		Use(&PetUseCase).
		Use(&Validator).
		Start(HandleRequest)
}
