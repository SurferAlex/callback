/** Supported client platforms for the install grid. */
export type PlatformId = "ios" | "android" | "macos" | "windows";

/** A downloadable client platform card. */
export interface Platform {
  id: PlatformId;
  /** Display name, e.g. "iOS". */
  name: string;
  /** Short marketing description. */
  description: string;
  /** Download / install URL. */
  url: string;
}
