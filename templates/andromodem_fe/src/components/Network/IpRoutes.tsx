import type { FC } from "react"

import type { IpRoute } from "@/types/network"

interface IpRoutesProps {
  ipRoutes: IpRoute[]
}

const IpRoutes: FC<IpRoutesProps> = ({ ipRoutes }) => {
  return (
    <div className="card bg-base-100 w-full shadow-sm">
      <div className="card-body">
        <h2 className="card-title">Ip Routes</h2>
        <span className="text-neutral-focus">
          IP Routes is the list of IP addresses that the device has.
        </span>
        {
          ipRoutes?.length > 0 ?
            <>
              <div className="overflow-x-auto rounded-box border border-base-content/5 bg-base-100">
                <table className="table">
                  <thead>
                    <tr>
                      <th>Interface</th>
                      <th>IP Address</th>
                    </tr>
                  </thead>
                  <tbody>
                    {
                      ipRoutes?.map((ipRoute) => (
                        <tr key={ipRoute.interface}>
                          <th>{ipRoute.interface}</th>
                          <td>{ipRoute.ip}</td>
                        </tr>
                      ))
                    }
                  </tbody>
                </table>
              </div>
            </> :
            <div className="text-center">
              <div className="text-neutral-focus">
                No IP Routes
              </div>
            </div>
        }
      </div>
    </div>
  )
}

export default IpRoutes