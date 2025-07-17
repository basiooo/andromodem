import { type FC } from "react"
import { LuRefreshCw } from "react-icons/lu"
import { MdBlock } from 'react-icons/md'
import { TbMessages } from 'react-icons/tb'

import useMessages from "@/hooks/useMessages"
import { type Device, DeviceState } from "@/types/device"

import MessageItemSkeleton from "@/components/Skeleton/MessageItemSkeleteon"
import MessageItem from "@/components/Messages/MessageItem"


type props = {
    device: Device
}

const MessageList: FC<props> = ({ device }) => {
    const { 
        messages, 
        finalMessages, 
        showValue, 
        showOptions, 
        isLoading, 
        isValidating, 
        error, 
        mutate, 
        handleShowOptionChange 
    } = useMessages(device)

    return (
        <div className="w-auto">
            <>
                <button onClick={() => mutate()} disabled={isLoading || isValidating || device.state === DeviceState.DISCONNECT} className="btn btn-xs md:btn-sm btn-active btn-primary mb-3">
                    <LuRefreshCw className={isLoading || isValidating ? "animate-spin" : ""} />
                    Refresh
                </button>
                {
                    isLoading || isValidating ?
                        <MessageItemSkeleton count={5} />
                        :
                        messages.length ?
                            <>
                                <div className="flex justify-between align-middle items-end">
                                    <span className="text-base md:text-lg font-bold mx-2">{messages.length <= 10 ? messages.length : showValue === -1 ? messages.length : showValue}/{messages.length}</span>
                                    <select defaultValue={messages.length <= 10 ? "*" : showValue} disabled={device.state !== DeviceState.ONLINE} className="select select-primary select-xs md:select-sm w-auto max-w-xs" onChange={(v: React.ChangeEvent<HTMLSelectElement>) => handleShowOptionChange(v.target.value)}>
                                        {
                                            messages.length <= 10 ?
                                                <option
                                                    key="*"
                                                    value="*">
                                                    Show All
                                                </option>
                                                :
                                                showOptions.map(value => {
                                                    if (messages.length >= value) {
                                                        return <option
                                                            key={value}
                                                            value={value} >{
                                                            value === -1 ? "Show All" : value
                                                        }</option>
                                                    }
                                                })
                                        }
                                    </select>
                                </div>
                                <div className="overflow-scroll h-96 mt-3">
                                    {
                                        finalMessages.map((message) => (
                                            <MessageItem key={`${device.serial}_${message.row}`} message={message} />
                                        ))
                                    }
                                </div>
                            </>
                            :
                            <div className='text-center'>
                                {
                                    error ?
                                        <>
                                            <MdBlock className='m-auto text-6xl md:text-8xl' color='red' />
                                            <span className='text-md md:text-xl'>{error.response.data.message}</span>
                                        </>
                                        :
                                        <>
                                            <TbMessages className='m-auto' fontSize={140} />
                                            <span className='text-xl'>No messages yet</span>
                                        </>
                                }
                            </div>
                }
            </>
        </div>
    )
}

export default MessageList