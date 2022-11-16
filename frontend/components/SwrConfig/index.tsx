import React from 'react';
import { SWRConfig } from 'swr';
import { apiCall, setUserAccessToken } from '../../helpers';

type Props = {
  children: React.ReactNode,
};

async function fetcher(endpoint: string, init?: RequestInit) {
  let response: Response;

  try {
    response = await apiCall(endpoint, init);
  } catch (e) {
    // TODO: handle better
    throw new Error('An error occurred while fetching the data.');
    // throw new SwrError('An error occurred while fetching the data.');
  }

  if (!response.ok) {
    if (response.status === 401) {
      setUserAccessToken(null);
    }

    // TODO: handle better
    const responseJson = await response.json();
    throw new Error('An error occurred while fetching the data.');
    // throw new SwrError('An error occurred while fetching the data.', response.status, responseJson);
  }

  return response.json();
}

export function SwrConfig({ children }: Props) {
  return (
    <SWRConfig value={{ fetcher, refreshInterval: 60000 }}>
      {children}
    </SWRConfig>
  );
}
