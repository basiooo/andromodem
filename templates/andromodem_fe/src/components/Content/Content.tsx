import { type FC } from "react"
import { GiNetworkBars } from "react-icons/gi"
import {MdMessage, MdPermDeviceInformation, MdSettingsPower} from "react-icons/md"
import { TbAutomation } from "react-icons/tb"

import DeviceInfo from "@/components/DeviceInfo/DeviceInfo"
import FeatureAvailability from "@/components/FeatureAvailability/FeatureAvailability"
import MessageList from "@/components/Messages/MessageList"
import Monitoring from "@/components/Monitoring/Monitoring"
import Network from "@/components/Network/Network"
import Power from "@/components/Power/Power"
import useTabManagement from "@/hooks/useTabManagement"
import { useDevicesStore } from "@/stores/devicesStore"
import { DeviceState } from "@/types/device"

const Content: FC = () => {
    const { deviceUsed } = useDevicesStore()
    const { activeTab, handleTabChange } = useTabManagement("device_info")

    if (deviceUsed == null){
        return (
            <div className="card card-compact w-full bg-base-200 shadow-xl mt-7 md:mt-10">
                <div className="card-body text-center">
                    <h1 className="text-md sm:text-lg md:text-xl font-bold">No Device Selected</h1>
                </div>
            </div>
        )
    }

    const tabs = [
        {
            key: "device_info",
            label: "Device Info",
            icon: <MdPermDeviceInformation className="w-4 h-4 mr-2" />
        },
        {
            key: "network",
            label: "Network",
            icon: <GiNetworkBars className="w-4 h-4 mr-2" />
        },
        {
            key: "messages",
            label: "Messages",
            icon: <MdMessage className="w-4 h-4 mr-2" />
        },
        {
            key: "monitoring",
            label: "Monitoring",
            icon: <TbAutomation className="w-4 h-4 mr-2" />
        },
        {
            key: "device_power",
            label: "Power",
            icon: <MdSettingsPower className="w-4 h-4 mr-2" />
        }
    ]

    const renderActiveTabContent = () => {
        switch (activeTab) {
            case "device_info":
                return <DeviceInfo device={deviceUsed}/>
            case "network":
                return <Network device={deviceUsed}/>
            case "messages":
                return <MessageList device={deviceUsed}/>
            case "monitoring":
                return <Monitoring device={deviceUsed}/>
            case "device_power":
                return <Power device={deviceUsed}/>
            default:
                return <DeviceInfo device={deviceUsed}/>
        }
    }

    return (
        <div className="w-full mt-5">
            {deviceUsed.state === DeviceState.DISCONNECT ? (
            <div role="alert" className="alert alert-error mb-4">
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="h-6 w-6 shrink-0 stroke-current"
                    fill="none"
                    viewBox="0 0 24 24"
                >
                    <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth="2"
                        d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                </svg>
                <span>Device Disconnected.</span>
            </div>
            ) : (
            <></>
            )}
            <FeatureAvailability device={deviceUsed} />
            <div role="tablist" className="tabs tabs-lift tabs-boxed mt-5">
                {tabs.map((tab) => (
                    <button
                        disabled={deviceUsed.state !== DeviceState.ONLINE}
                        key={tab.label}
                        role="tab"
                        className={`tab flex items-center ${activeTab === tab.key ? "tab-active !bg-base-200 text-lg" : ""}`}
                        onClick={() => handleTabChange(tab.key)}
                    >
                        <span className="md:hidden">{tab.icon}</span>
                        <span className="hidden md:block">{tab.label}</span>
                    </button>
                ))}
            </div>

            <div className="card card-compact w-full bg-base-200 shadow-xl">
                <div className="card-body">
                    <h2 className="card-title text-sm md:text-md mb-2 border-b-3 border-base-100">{deviceUsed.model}</h2>
                    {renderActiveTabContent()}
                </div>
            </div>
        </div>
    )
}

export default Content