import Axios from "axios";
import { backendBaseUrl } from "shared/config";

export type Container = {
  ContainerIp: string,
  PingTimeMs: number,
  LastSuccessfulPing: string | null,
  Status: string,
}

export async function getContainersList(): Promise<Container[]> {
  const { data } = await Axios.get<Container[]>(backendBaseUrl + "/containers");
  return data
}
