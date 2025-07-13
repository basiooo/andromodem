import { Anchorme } from "react-anchorme"

import type { Message } from "@/types/message"


type props = {
    message: Message
}

const CustomLink = (props: object) => {
    return (
        <a className="link" {...props} />
    )
}
const MessageItem = ({ message }: props) => {
    return (
        <div className="collapse collapse-arrow  bg-base-300 my-3">
            <input type="checkbox" className="peer" />
            <div className="collapse-title text-md md:text-lg bg-base-300 peer-checked:bg-base-300">
                {message.address}
                <p className='text-xs md:text-sm text-gray-500'>
                    {message.date}
                </p>
            </div>
            <div className="collapse-content bg-base-300 border-t pt-3 md:text-lg" style={{ wordBreak: 'break-word' }}>
                <Anchorme linkComponent={CustomLink} target="_blank" rel="noreferrer noopener">
                    {message.body}
                </Anchorme>
            </div>
        </div>
    )
}

export default MessageItem