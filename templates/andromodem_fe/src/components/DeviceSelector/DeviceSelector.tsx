import type { FC } from "react"
import { useState } from "react"

import useListenDevices from "@/hooks/useListenDevices"
import { useDevicesStore } from "@/stores/devicesStore"
import { DeviceState } from "@/types/device"

const DeviceSelector: FC = () => {
    const { devices, isConnected } = useListenDevices()
    const { setDeviceUsed, unsetDeviceUsed } = useDevicesStore()

    const [selectedDevice, setSelectedDevice] = useState<string>('DEFAULT')

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
            </div>
        </div>
    )
}
export default DeviceSelector
