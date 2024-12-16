import { useState, useEffect } from 'react';
import './App.css';

function App() {
  const [searchTerm, setSearchTerm] = useState('');
  const [beers, setBeers] = useState([]);
  const [loading, setLoading] = useState(false);

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

  return (
    <div className="App">
      <h1>Beer Finder</h1>
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
    </div>
  );
}

export default App;
