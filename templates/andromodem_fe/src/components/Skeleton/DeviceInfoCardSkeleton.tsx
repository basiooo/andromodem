import type { CSSProperties } from "react"

type Props = {
  count?: number;
};

const DeviceInfoCardSkeleton = ({ count = 4 }: Props) => {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5 md:gap-2 justify-items-center">
      {Array.from({ length: count }).map((_, i) => (
        <div key={i} className="card bg-base-100 w-full md:w-auto shadow-sm animate-pulse">
          <div className="card-body">
            <h2 className="card-title">
              <div className="h-4 w-24 bg-base-300 rounded" />
            </h2>
            <div className="flex items-center justify-around gap-4">
              <div className="bg-base-300 rounded-full p-2">
                <div
                  className="radial-progress text-base-300"
                  style={
                    {
                      "--value": 50,
                      "--size": "4rem"
                    } as CSSProperties
                  }
                >
                  <span className="invisible">%</span>
                </div>
              </div>
              <div className="md:text-xs space-y-1">
                <div className="h-3 w-24 bg-base-300 rounded" />
                <div className="h-3 w-20 bg-base-300 rounded" />
              </div>
            </div>
            <div className="text-center mt-4">
              <div className="btn btn-sm bg-base-300 text-base-300 border-none w-20 h-6" />
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}

export default DeviceInfoCardSkeleton
