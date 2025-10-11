import React, { useState } from 'react';
import axios from 'axios';
import './App.css';

function App() {
  const [event, setEvent] = useState({ user_id: '', type: '', metadata: '' });
  const [query, setQuery] = useState({ user_id: '', start: '', end: '' });
  const [events, setEvents] = useState([]);

  const handleCreate = async () => {
    let metadataObj;
    try {
      metadataObj = JSON.parse(event.metadata || '{}'); // Parse string to object; default to empty object
    } catch (e) {
      alert('Invalid metadata JSON');
      return;
    }
    const eventToSend = { 
      user_id: parseInt(event.user_id, 10), // Ensure user_id is an integer
      type: event.action,
      metadata: metadataObj 
    };
    if (isNaN(eventToSend.user_id)) {
      alert('Invalid User ID: must be a number');
      return;
    }
    try {
      await axios.post('http://localhost:8080/events', eventToSend);
      alert('Event created');
    } catch (err) {
      alert('Error creating event: ' + (err.response?.data?.error || err.message));
    }
  };

  const handleRetrieve = async () => {
    try {
      const params = { ...query };
      if (query.user_id) params.user_id = parseInt(query.user_id, 10); // Ensure user_id is an integer
      const res = await axios.get('http://localhost:8080/events', { params });
      setEvents(res.data);
    } catch (err) {
      alert('Error retrieving events: ' + (err.response?.data?.error || err.message));
    }
  };

  return (
    <div className="App">
      <h1>User Activity Tracker</h1>
      
      <h2>Create Event</h2>
      <input 
        type="number" 
        placeholder="User ID" 
        value={event.user_id} 
        onChange={e => setEvent({...event, user_id: e.target.value})} 
      />
      <input 
        placeholder="Action" 
        value={event.action} 
        onChange={e => setEvent({...event, type: e.target.value})} 
      />
      <input 
        placeholder='Metadata (JSON, e.g. {"page":"/home"})' 
        value={event.metadata} 
        onChange={e => setEvent({...event, metadata: e.target.value})} 
      />
      <button onClick={handleCreate}>Create</button>
      
      <h2>Retrieve Events</h2>
      <input 
        type="number" 
        placeholder="User ID" 
        value={query.user_id} 
        onChange={e => setQuery({...query, user_id: e.target.value})} 
      />
      <input 
        placeholder="Start (ISO, e.g. 2025-10-01T00:00:00Z)" 
        value={query.start} 
        onChange={e => setQuery({...query, start: e.target.value})} 
      />
      <input 
        placeholder="End (ISO, e.g. 2025-10-10T23:59:59Z)" 
        value={query.end} 
        onChange={e => setQuery({...query, end: e.target.value})} 
      />
      <button onClick={handleRetrieve}>Retrieve</button>
      
      <ul>
        {events.map((ev, i) => (
          <li key={i}>{JSON.stringify(ev)}</li>
        ))}
      </ul>
    </div>
  );
}

export default App;