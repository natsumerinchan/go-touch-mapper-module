package main

// https://github.com/kp7742/TouchSimulation/blob/main/UinputDefs.go

import (
	"syscall"
)

//---------------------------------EVCodes--------------------------------------//

//Ref: input-event-codes.h
const (
	evSyn           = 0x00
	evKey           = 0x01
	evRel           = 0x02
	evAbs           = 0x03
	evMsc           = 0x04
	relX            = 0x00
	relY            = 0x01
	relWheel        = 0x08
	relHWheel       = 0x06
	btnTouch        = 0x14a
	synReport       = 0
	synMtReport     = 2
	absMtSlot       = 0x2f
	absMtPositionX  = 0x35
	absMtPositionY  = 0x36
	absMtTrackingId = 0x39
	absMax          = 0x3f
	absCnt          = absMax + 1
	inputPropDirect = 0x01
	inputPropMax    = 0x1f
	inputPropCnt    = inputPropMax + 1
	//-------------------------------appends-------------------------------//
	absMtTouchMajor = 0x30
	absMtWidthMajor = 0x32
)

//---------------------------------IOCTL--------------------------------------//

//Ref: ioctl.h
const (
	iocNone  = 0x0
	iocWrite = 0x1
	iocRead  = 0x2

	iocNrbits   = 8
	iocTypebits = 8
	iocSizebits = 14
	iocNrshift  = 0

	iocTypeshift = iocNrshift + iocNrbits
	iocSizeshift = iocTypeshift + iocTypebits
	iocDirshift  = iocSizeshift + iocSizebits
)

func _IOC(dir int, t int, nr int, size int) int {
	return (dir << iocDirshift) | (t << iocTypeshift) |
		(nr << iocNrshift) | (size << iocSizeshift)
}

func _IOR(t int, nr int, size int) int {
	return _IOC(iocRead, t, nr, size)
}

func _IOW(t int, nr int, size int) int {
	return _IOC(iocWrite, t, nr, size)
}

//Ref: input.h
func EVIOCGNAME() int {
	return _IOC(iocRead, 'E', 0x06, uinputMaxNameSize)
}

func EVIOCGPROP() int {
	return _IOC(iocRead, 'E', 0x09, inputPropMax)
}

func EVIOCGABS(abs int) int {
	return _IOR('E', 0x40+abs, 24) //sizeof(struct input_absinfo)
}

func EVIOCGBIT(ev, len int) int {
	return _IOC(iocRead, 'E', 0x20+ev, len)
}

func EVIOCGRAB() int {
	return _IOW('E', 0x90, 4) //sizeof(int)
}

//Syscall
func ioctl(fd uintptr, name int, data uintptr) error {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(name), data)
	if err != 0 {
		return err
	}
	return nil
}

//---------------------------------Input--------------------------------------//

type InputID struct {
	BusType uint16
	Vendor  uint16
	Product uint16
	Version uint16
}

type AbsInfo struct {
	Value      int32
	Minimum    int32
	Maximum    int32
	Fuzz       int32
	Flat       int32
	Resolution int32
}

type InputEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

//---------------------------------UInput--------------------------------------//

//Ref: uinput.h
const (
	uinputMaxNameSize = 80
)

type UinputUserDev struct {
	Name       [uinputMaxNameSize]byte
	ID         InputID
	EffectsMax uint32
	AbsMax     [absCnt]int32
	AbsMin     [absCnt]int32
	AbsFuzz    [absCnt]int32
	AbsFlat    [absCnt]int32
}

//Ref: uinput.h
func UISETEVBIT() int {
	return _IOW('U', 100, 4) //sizeof(int)
}

func UISETKEYBIT() int {
	return _IOW('U', 101, 4) //sizeof(int)
}

func UISETEVRELBIT() int {
	return _IOW('U', 102, 4) //sizeof(int)
}

func UISETABSBIT() int {
	return _IOW('U', 103, 4) //sizeof(int)
}

func UISETPROPBIT() int {
	return _IOW('U', 110, 4) //sizeof(int)
}

func UIDEVCREATE() int {
	return _IOC(iocNone, 'U', 1, 0)
}

func UIDEVDESTROY() int {
	return _IOC(iocNone, 'U', 2, 0)
}
