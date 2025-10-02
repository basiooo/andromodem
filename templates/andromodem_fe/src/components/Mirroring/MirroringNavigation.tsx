import { type FC } from "react"
import { AiOutlineFullscreenExit } from "react-icons/ai"
import { IoChevronBackOutline } from "react-icons/io5"
import { MdPowerSettingsNew } from "react-icons/md"
import { RiHomeLine } from "react-icons/ri"
import { TbBoxMultiple } from "react-icons/tb"

import { KeyCommand, type KeyCommandValue } from "@/types/mirroring"

interface MirroringNavigationProps {
    isConnected: boolean
    isFullscreen: boolean
    onToggleFullscreen: () => void
    onSendKeyCommand?: (key: KeyCommandValue) => void
}

const MirroringNavigation: FC<MirroringNavigationProps> = ({
    isConnected,
    isFullscreen,
    onToggleFullscreen,
    onSendKeyCommand
}) => {
    const handleKeyPress = (key: KeyCommandValue) => {
        if (isConnected && onSendKeyCommand) {
            onSendKeyCommand(key)
        }
    }
    return (
        <div className="flex justify-center items-center space-x-4 py-3 rounded-md bg-opacity-90 backdrop-blur-sm">
            <button
                onClick={() => handleKeyPress(KeyCommand.RECENT)}
                disabled={!isConnected}
                className="flex items-center justify-center w-8 h-8 rounded-full bg-gray-700 hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                title="Recent Apps"
            >
                <TbBoxMultiple className="w-6 h-6 text-white" />
            </button>
            <button
                onClick={() => handleKeyPress(KeyCommand.HOME)}
                disabled={!isConnected}
                className="flex items-center justify-center w-8 h-8 rounded-full bg-gray-700 hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                title="Home"
            >
                <RiHomeLine className="w-6 h-6 text-white" />
            </button>
            <button
                onClick={() => handleKeyPress(KeyCommand.BACK)}
                disabled={!isConnected}
                className="flex items-center justify-center w-8 h-8 rounded-full bg-gray-700 hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                title="Back"
            >
                <IoChevronBackOutline className="w-6 h-6 text-white" />
            </button>
            <button
                onClick={() => handleKeyPress(KeyCommand.POWER)}
                disabled={!isConnected}
                className="flex items-center justify-center w-8 h-8 rounded-full bg-red-700 hover:bg-red-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                title="Power"
            >
                <MdPowerSettingsNew className="w-6 h-6 text-white" />
            </button>
            {
                isFullscreen && (
                    <button
                        onClick={() => onToggleFullscreen()}
                        disabled={!isConnected}
                        className="flex items-center justify-center w-8 h-8 rounded-full bg-gray-700 hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                        title="Exit Fullscreen"
                    >
                        <AiOutlineFullscreenExit className="w-6 h-6 text-white" />
                    </button>
                )
            }
        </div>
    )
}
export default MirroringNavigation