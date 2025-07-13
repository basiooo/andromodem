type props = {
    count: number
}
const MessageItemSkeleton= ({ count = 2 }: props) => {
    const skeletons = Array.from({ length: count }, (_, index) => (
        <div className="skeleton w-full min-h-20 p-5 my-2" key={index}>
        </div>
    ))

    return <>{skeletons}</>
}

export  default  MessageItemSkeleton