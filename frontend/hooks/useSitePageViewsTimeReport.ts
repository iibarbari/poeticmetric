import { useMemo } from "react";
import useSWR, { SWRResponse } from "swr";
import { hydrateSitePageViewsTimeReport } from "../helpers";
import { useReportQueryParams } from "./useReportQueryParams";

type HydratedSwrResponse = SWRResponse<SitePageViewsTimeReport> & {
  hydratedData?: HydratedSitePageViewsTimeReport,
};

export function useSitePageViewsTimeReport(): HydratedSwrResponse {
  const reportQueryParams = useReportQueryParams();
  const swrResponse = useSWR<SitePageViewsTimeReport>(`/site-reports/page-views-time?${reportQueryParams}`);

  const hydratedData = useMemo<HydratedSitePageViewsTimeReport | undefined>(() => {
    if (swrResponse.data === undefined) {
      return undefined;
    }

    return hydrateSitePageViewsTimeReport(swrResponse.data);
  }, [swrResponse.data]);

  return { ...swrResponse, hydratedData };
}
