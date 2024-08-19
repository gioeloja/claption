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

  // call AWS backend
  const apiResponse = await fetch(`${apiUrl}/${jobID}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  });

  const apiData = await apiResponse.json();

  return NextResponse.json({ apiData });
}
