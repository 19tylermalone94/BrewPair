import { useState, useEffect } from 'react';

function App() {
  const [searchTerm, setSearchTerm] = useState('');
  const [beers, setBeers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [selectedImage, setSelectedImage] = useState<File | null>(null);
  const [identifiedBeer, setIdentifiedBeer] = useState<any | null>(null); // Adjusted type to handle beer object

  useEffect(() => {
    if (searchTerm.trim() === '') {
      setBeers([]);
      return;
    }

    const fetchBeers = async () => {
      setLoading(true);
      try {
        const response = await fetch(`http://localhost:8080/beers?search=${searchTerm}`);
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        const data = await response.json();
        setBeers(data);
      } catch (error) {
        console.error('Error fetching beers:', error);
        setBeers([]);
      } finally {
        setLoading(false);
      }
    };

    fetchBeers();
  }, [searchTerm]);

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    setSelectedImage(file || null);
  };

  const handleImageUpload = async () => {
    if (!selectedImage) {
      alert('Please select an image before uploading.');
      return;
    }

    const formData = new FormData();
    formData.append('image', selectedImage);

    try {
      const response = await fetch('http://localhost:8080/identify-beer', {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) {
        const errorText = await response.text(); // Log the error response
        throw new Error(`Failed to identify beer. Server response: ${errorText}`);
      }

      const data = await response.json();
      console.log('API Response:', data); // Debug: log the full API response

      if (data && data.name) {
        setIdentifiedBeer(data); // Save the full beer object
      } else {
        setIdentifiedBeer({ name: 'Unknown' }); // Handle cases where no name is returned
      }
    } catch (error) {
      console.error('Error identifying beer:', error);
      setIdentifiedBeer({ name: 'Error occurred' });
    }
  };

  return (
    <div className="App">
      <h1>Brew Pair</h1>

      {/* Search for beers */}
      <div className="search-container">
        <input
          type="text"
          placeholder="Search for beers..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="search-box"
        />
      </div>

      {loading ? (
        <p>Loading...</p>
      ) : (
        <div className="beer-list">
          <ul>
            {beers.map((beer: any) => (
              <li key={beer.id}>
                <strong>{beer.name}</strong> - {beer.style}
                {beer.ibu && <span> (IBU: {beer.ibu})</span>}
              </li>
            ))}
          </ul>
        </div>
      )}

      {/* Upload image to identify beer */}
      <div className="upload-container">
        <h2>Identify a Beer</h2>
        <input type="file" accept="image/*" onChange={handleImageChange} />
        <button onClick={handleImageUpload}>Upload and Identify</button>
        {identifiedBeer && (
          <div className="identified-beer">
            <h3>Identified Beer:</h3>
            {identifiedBeer.name ? (
              <p>
                <strong>{identifiedBeer.name}</strong>
                {identifiedBeer.style && ` - ${identifiedBeer.style}`}
                {identifiedBeer.ibu && ` (IBU: ${identifiedBeer.ibu})`}
              </p>
            ) : (
              <p>Unknown</p>
            )}
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
