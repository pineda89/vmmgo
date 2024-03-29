package vmmgo

import "unsafe"

func (inst *VMM) MapVad(pid uint32) (*VMMDLL_MAP_VAD, error) {
	data := make([]byte, 0x4000000)

	result, _, _ := call("VMMDLL_Map_GetVadU", initializeResult, uintptr(pid), 1, uintptr(unsafe.Pointer(&data)))
	if result == 0 {
		return nil, ERR_CALL
	}

	oneElement := (*vMMDLL_MAP_VAD)(unsafe.Pointer(&data[0]))

	if oneElement.DwVersion != VMMDLL_MAP_VAD_VERSION {
		return nil, ERR_BAD_VERSION
	}

	sizeofchild := unsafe.Sizeof(VMMDLL_MAP_VADENTRY{})

	returnStruct := &VMMDLL_MAP_VAD{
		DwVersion:   oneElement.DwVersion,
		CMap:        oneElement.CMap,
		CbMultiText: oneElement.CbMultiText,
		_Reserved1:  oneElement._Reserved1,
		CPage:       oneElement.CPage,
		PbMultiText: oneElement.PbMultiText,
		PMap:        make([]VMMDLL_MAP_VADENTRY, oneElement.CMap),
	}

	addr := uintptr(unsafe.Pointer(&oneElement.PMap[0]))
	for i := 0; i < int(oneElement.CMap); i++ {
		from := addr + (uintptr(i) * sizeofchild)

		t := (*VMMDLL_MAP_VADENTRY)(unsafe.Pointer(from))

		t.UszText = readString(unsafe.Pointer(t.UszText))

		returnStruct.PMap[i] = *t
	}

	return returnStruct, nil
}
