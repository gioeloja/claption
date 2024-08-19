import { NextResponse } from 'next/server';

export async function GET(request: Request, { params }: { params: { jobID: string } }) {
  const { jobID } = params;
  const apiUrl = process.env.STATUS_POLL_API;

  if (!apiUrl) {
    return new Response(JSON.stringify({ error: 'API URL not configured in .env' }), {
      status: 500,
      headers: { 'Content-Type': 'application/json' },
    });
  }

  let attempts = 0;
  let maxAttempts = 2;
  let apiData;
  let apiResponse;
  let contentType;

  while (attempts < maxAttempts) {
    apiResponse = await fetch(`${apiUrl}/${jobID}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    contentType = apiResponse.headers.get('content-type');

    if (contentType && contentType.includes('application/json')) {
      apiData = await apiResponse.json();
      break;
    } else {
      apiData = await apiResponse.text();
    }

    attempts += 1;

    if (attempts < maxAttempts && (!contentType || !contentType.includes('application/json'))) {
      // wait between retries
      await new Promise(resolve => setTimeout(resolve, 1000));
    }
  }

  if (attempts === maxAttempts && (!contentType || !contentType.includes('application/json'))) {
    return new Response(JSON.stringify({ error: 'Failed to retrieve JSON response after 2 attempts' }), {
      status: 500,
      headers: { 'Content-Type': 'application/json' },
    });
  }

  return NextResponse.json({ apiData });
}
