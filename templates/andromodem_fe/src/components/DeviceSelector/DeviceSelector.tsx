import type { FC } from "react"
import { useEffect,useState } from "react"

import useListenDevices from "@/hooks/useListenDevices"
import { useDevicesStore } from "@/stores/devicesStore"
import { DeviceState } from "@/types/device"

const DeviceSelector: FC = () => {
    const { devices, isConnected } = useListenDevices()
    const { setDeviceUsed, deviceUsed, unsetDeviceUsed } = useDevicesStore()

    const [selectedDevice, setSelectedDevice] = useState<string>('DEFAULT')
    const [selectedDeviceSerial, setSelectedDeviceSerial] = useState<string | null>(null)

    const autoSelectEnabled = selectedDeviceSerial === deviceUsed?.serial

    useEffect(() => {
        const savedDeviceSerial = localStorage.getItem('selectedDeviceSerial')
        setSelectedDeviceSerial(savedDeviceSerial)
        if (savedDeviceSerial && savedDeviceSerial !== 'DEFAULT') {
            setSelectedDevice(savedDeviceSerial)
        }
    }, [])

    useEffect(() => {
        const savedDeviceSerial = localStorage.getItem('selectedDeviceSerial')
        setSelectedDeviceSerial(savedDeviceSerial)
    }, [selectedDevice, deviceUsed])

    useEffect(() => {
        if (devices.length > 0 && isConnected) {
            const savedDeviceSerial = localStorage.getItem('selectedDeviceSerial')

            if (savedDeviceSerial && savedDeviceSerial !== 'DEFAULT') {
                const savedDevice = devices.find(device =>
                    device.serial === savedDeviceSerial &&
                    device.state === DeviceState.ONLINE
                )

                if (savedDevice) {
                    setSelectedDevice(savedDeviceSerial)
                    setDeviceUsed(savedDevice)
                    setSelectedDeviceSerial(savedDeviceSerial)
                }
            }
        }
    }, [devices, isConnected, setDeviceUsed])

    const handleDeviceSelectChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
        const serial = e.target.value
        setSelectedDevice(serial)

        const device = devices.find((device) => device.serial === serial)
        if (device) {
            setDeviceUsed(device)
        } else {
            unsetDeviceUsed()
        }
    }

    const handleAutoSelectChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.checked) {
            localStorage.setItem('selectedDeviceSerial', selectedDevice)
        } else {
            localStorage.removeItem('selectedDeviceSerial')
        }
        const savedDeviceSerial = localStorage.getItem('selectedDeviceSerial')
        setSelectedDeviceSerial(savedDeviceSerial)
    }

    return (
        <div className="card card-compact w-full bg-base-200 shadow-xl">
            <div className="card-body">
                <h2 className="card-title text-md md:text-xl">Select Device
                    {isConnected ? (
                        <div className="tooltip" data-tip="Server Connected">
                            <div aria-label="success" className="status status-success"></div>
                        </div>
                    ) : (
                        <div className="tooltip" data-tip="Server Not Connected">
                            <div aria-label="error" className="status status-error"></div>
                        </div>
                    )}
                </h2>
                <select
                    className="select select-primary w-full sm:select-lg"
                    value={selectedDevice}
                    onChange={handleDeviceSelectChange}
                >
                    <option value="DEFAULT"
                        id="device_selector">{`Select device. ${devices.length} device detected`}</option>
                    {devices.map((device) => (
                        <option
                            key={device.serial}
                            value={device.serial}
                            disabled={device.state !== DeviceState.ONLINE || !isConnected}
                        >
                            {device.model ? `${device.model} / ${device.serial}` : `${device.serial}`}
                            {`(${device.state})`}
                        </option>
                    ))}
                </select>
                {
                    deviceUsed ?
                        <div className="form-control">
                            <label className="label cursor-pointer justify-start gap-2">
                                <input
                                    disabled={deviceUsed.state !== DeviceState.ONLINE}
                                    type="checkbox"
                                    className="checkbox checkbox-primary checkbox-xs md:checkbox-sm"
                                    checked={autoSelectEnabled}
                                    onChange={handleAutoSelectChange}
                                />
                                <span className="text-xs md:text-md">Auto select when open AndroModem</span>
                            </label>
                        </div>
                        : <></>
                }
            </div>
        </div>
    )
}
export default DeviceSelector
