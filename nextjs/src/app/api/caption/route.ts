export async function POST(req: Request) {
    try {
      const { base64encoding } = await req.json();
  
      const apiUrl = process.env.CAPTION_IMAGE_API;
  
      if (!apiUrl) {
        return new Response(JSON.stringify({ error: 'API URL not configured in .env' }), {
          status: 500,
          headers: { 'Content-Type': 'application/json' },
        });
      }
  
      // call AWS backend
      const apiResponse = await fetch(apiUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ body: base64encoding }),
      });
  
      if (!apiResponse.ok) {
        const apiError = await apiResponse.text();
        console.error("Error from API:", apiError);
        return new Response(JSON.stringify({ error: 'Error from external API' }), {
          status: apiResponse.status,
          headers: { 'Content-Type': 'application/json' },
        });
      }
  
      const apiData = await apiResponse.json();

      return new Response(JSON.stringify({ apiData }), {
        headers: { 'Content-Type': 'application/json' },
      });
    } catch (error) {
      console.error("Error processing request:", error);
      return new Response(JSON.stringify({ error: 'Internal Server Error' }), {
        status: 500,
        headers: { 'Content-Type': 'application/json' },
      });
    }
  }
  