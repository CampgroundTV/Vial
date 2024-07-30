        async function fetchStreams() {
            try {
                const response = await fetch('http://localhost:8080/streams');
                const streams = await response.json();
                const streamsDiv = document.getElementById('streams');
                streamsDiv.innerHTML = '';
                streams.forEach(stream => {
                    const streamDiv = document.createElement('div');
                    streamDiv.className = 'stream';
                    streamDiv.innerHTML = `<strong>${stream.title}</strong> (${stream.category}) - ${stream.is_live ? 'Live' : 'Offline'}`;
                    streamsDiv.appendChild(streamDiv);
                });
            } catch (error) {
                console.error('Error fetching streams:', error);
            }
        }

        setInterval(fetchStreams, 30000); // Fetch streams every 30 seconds
        fetchStreams(); // Initial fetch
