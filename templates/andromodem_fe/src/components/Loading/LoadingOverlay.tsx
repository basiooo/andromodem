import type { FC } from "react"

const LoadingOverlay: FC = () => {
    return (
        <div
            className="absolute inset-0 flex  bg-base-300/40
    items-center justify-center
    text-white z-20"
        >
            <div
                className="bg-gray-800 p-8
        rounded-lg shadow-lg"
            >
                <span className="loading loading-bars loading-xl"></span>
            </div>
        </div>
    )
}

export default LoadingOverlay