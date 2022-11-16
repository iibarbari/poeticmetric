import dayjs from "dayjs";
import duration from "dayjs/plugin/duration";
import isoWeek from "dayjs/plugin/isoWeek";
import localizedFormat from "dayjs/plugin/localizedFormat";
import relativeTime from "dayjs/plugin/relativeTime";
import { AppProps } from "next/app";
import React, { useEffect } from "react";
import { AuthAndApiHandler, SwrConfig, ToastsHandler } from "../components";
import "../styles/style.scss";

export default function App({ Component, pageProps }: AppProps) {
  useEffect(() => {
    dayjs.extend(duration);
    dayjs.extend(isoWeek);
    dayjs.extend(localizedFormat);
    dayjs.extend(relativeTime);
  }, []);

  return (
    <SwrConfig>
      <AuthAndApiHandler>
        <ToastsHandler>
          <Component {...pageProps} />
        </ToastsHandler>
      </AuthAndApiHandler>
    </SwrConfig>
  );
}
