package device

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TODO: Add mockery to the project
// MockDevice is a testify/mock implementation of the Device interface.
type MockDevice struct {
	mock.Mock
}

func (m *MockDevice) Read(p []byte) (int, error) {
	args := m.Called(p)
	// Copy provided byte slice into p if present
	if data, ok := args.Get(0).([]byte); ok {
		copy(p, data)
		return len(data), args.Error(1)
	}
	return 0, args.Error(1)
}

func (m *MockDevice) Close() error {
	args := m.Called()
	return args.Error(0)
}

// helper: creates a valid 8-byte signal buffer
// Signal format according to model.Signal: Byte[0]=Type, Byte[1]=High, Byte[2]=Low, ...
// Value = (High << 8) | Low
func makeSignalBuffer(signalType byte, value uint16) []byte {
	high := byte(value >> 8)
	low := byte(value & 0xFF)
	return []byte{signalType, high, low, 0x0D, 0, 0, 0, 0}
}

// --- Tests for Read ---
func TestRead_CarbonDioxide(t *testing.T) {
	// CO2 value: 800 ppm
	co2Value := uint16(800)
	buf := makeSignalBuffer(OperationCarbonDioxide, co2Value)

	d := new(MockDevice)
	d.On("Read", mock.AnythingOfType("[]uint8")).Return(buf, nil).Times(4)

	ac := AirController{device: d}

	co2, _, _, err := ac.Read(context.Background())
	require.NoError(t, err)
	assert.Equal(t, float64(co2Value), co2)
	d.AssertExpectations(t)
}

func TestRead_Temperature(t *testing.T) {
	// Temperature value: raw value 4619 → (4619/16) - 273.15 = 288.6875 - 273.15 = 15.54°C
	rawTemp := uint16(4619)
	expectedTemp := float64(rawTemp)/16 - 273.15
	buf := makeSignalBuffer(OperationTemperature, rawTemp)

	d := new(MockDevice)
	d.On("Read", mock.AnythingOfType("[]uint8")).Return(buf, nil).Times(4)

	ac := AirController{device: d}

	_, temp, _, err := ac.Read(context.Background())
	require.NoError(t, err)
	assert.Equal(t, expectedTemp, temp)
	d.AssertExpectations(t)
}

func TestRead_Humidity(t *testing.T) {
	// Humidity value: raw value 5500 → 5500/100 = 55.0 %
	rawHumidity := uint16(5500)
	expectedHumidity := float64(rawHumidity) / 100
	buf := makeSignalBuffer(OperationHumidity, rawHumidity)

	d := new(MockDevice)
	d.On("Read", mock.AnythingOfType("[]uint8")).Return(buf, nil).Times(4)

	ac := AirController{device: d}

	_, _, humidity, err := ac.Read(context.Background())
	require.NoError(t, err)
	assert.Equal(t, expectedHumidity, humidity)
	d.AssertExpectations(t)
}

func TestRead_AllThreeMeasurements(t *testing.T) {
	rawCO2 := uint16(950)
	rawTemp := uint16(4750) // (4750/16) - 273.15 ≈ 23.60°C
	rawHum := uint16(6000)  // 60.0 %
	unknownBuf := []byte{0xFF, 0x00, 0x00, 0x0D, 0, 0, 0, 0}

	d := new(MockDevice)
	// 4 reads: CO2, temp, humidity + one unknown type → continue
	d.On("Read", mock.AnythingOfType("[]uint8")).
		Return(makeSignalBuffer(OperationCarbonDioxide, rawCO2), nil).Once()
	d.On("Read", mock.AnythingOfType("[]uint8")).
		Return(makeSignalBuffer(OperationTemperature, rawTemp), nil).Once()
	d.On("Read", mock.AnythingOfType("[]uint8")).
		Return(makeSignalBuffer(OperationHumidity, rawHum), nil).Once()
	d.On("Read", mock.AnythingOfType("[]uint8")).
		Return(unknownBuf, nil).Once()

	ac := AirController{device: d}

	co2, temp, humidity, err := ac.Read(context.Background())
	require.NoError(t, err)
	assert.Equal(t, float64(rawCO2), co2, "CO2 mismatch")
	assert.Equal(t, float64(rawTemp)/16-273.15, temp, "temperature mismatch")
	assert.Equal(t, float64(rawHum)/100, humidity, "humidity mismatch")
	d.AssertExpectations(t)
}

func TestRead_DeviceReadError(t *testing.T) {
	readErr := errors.New("device not readable")

	d := new(MockDevice)
	d.On("Read", mock.AnythingOfType("[]uint8")).Return(nil, readErr).Once()

	ac := AirController{device: d}

	_, _, _, err := ac.Read(context.Background())
	require.Error(t, err)
	assert.ErrorIs(t, err, readErr)
	d.AssertExpectations(t)
}

func TestRead_ContextCancelledBeforeStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	d := new(MockDevice)
	// Read should never be called since the context is already cancelled
	ac := AirController{device: d}

	_, _, _, err := ac.Read(ctx)
	assert.ErrorIs(t, err, context.Canceled)
	d.AssertNotCalled(t, "Read", mock.Anything)
}

func TestRead_ContextDeadlineExceeded(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	time.Sleep(1 * time.Millisecond) // let the timeout expire

	d := new(MockDevice)
	ac := AirController{device: d}

	_, _, _, err := ac.Read(ctx)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	d.AssertNotCalled(t, "Read", mock.Anything)
}

func TestRead_UnknownSignalTypeIsSkipped(t *testing.T) {
	unknownBuf := []byte{0xAB, 0x00, 0x64, 0x0D, 0, 0, 0, 0}
	co2Buf := makeSignalBuffer(OperationCarbonDioxide, 500)

	d := new(MockDevice)
	// 4 reads: unknown, unknown, unknown, CO2
	d.On("Read", mock.AnythingOfType("[]uint8")).Return(unknownBuf, nil).Times(3)
	d.On("Read", mock.AnythingOfType("[]uint8")).Return(co2Buf, nil).Once()

	ac := AirController{device: d}

	co2, _, _, err := ac.Read(context.Background())
	require.NoError(t, err)
	assert.Equal(t, float64(500), co2)
	d.AssertExpectations(t)
}

// --- Tests for Close ---

func TestClose_Success(t *testing.T) {
	d := new(MockDevice)
	d.On("Close").Return(nil).Once()

	ac := AirController{device: d}

	err := ac.Close(context.Background())
	require.NoError(t, err)
	d.AssertExpectations(t)
}

func TestClose_DeviceError(t *testing.T) {
	closeErr := errors.New("device could not be closed")

	d := new(MockDevice)
	d.On("Close").Return(closeErr).Once()

	ac := AirController{device: d}

	err := ac.Close(context.Background())
	assert.ErrorIs(t, err, closeErr)
	d.AssertExpectations(t)
}

func TestClose_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	d := new(MockDevice)
	ac := AirController{device: d}

	err := ac.Close(ctx)
	assert.ErrorIs(t, err, context.Canceled)
	// Close() must not be called on a cancelled context
	d.AssertNotCalled(t, "Close")
}

func TestClose_ContextDeadlineExceeded(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	time.Sleep(1 * time.Millisecond) // let the timeout expire

	d := new(MockDevice)
	ac := AirController{device: d}

	err := ac.Close(ctx)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	d.AssertNotCalled(t, "Close")
}

// Ensure MockDevice satisfies the io.ReadCloser / Device interface at compile time.
var _ io.ReadCloser = (*MockDevice)(nil)
