import { useEffect, useState } from 'react'
import './App.css'
import { Movie } from './modules/Movie'

interface MovieInterface {
  id: Number
  title: string
}

function App() {
  const [count, setCount] = useState(0)
  const [searchData, setSearchData] = useState("")
  const [requestInputData, setRequestInputData] = useState("")
  const [movieArray, setMovieArray] = useState([])
  useEffect(() => {
    const timeOutId = setTimeout(() => setRequestInputData(searchData), 500)
    return () => clearTimeout(timeOutId)
  }, [searchData])

  useEffect(() => {
    if (requestInputData != ""){
      console.log("Tutaj leci request")
      console.log(requestInputData)
      const requestOptions = {
        method: 'POST',
        body: JSON.stringify({ searchbar_input: requestInputData })
      }
      fetch('http://127.0.0.1:8080/search', requestOptions)
          .then(response => response.json())
          .then(data => setMovieArray(data['titles']))
      }
  }, [requestInputData])

  const movies = movieArray.map((movie:MovieInterface) => <Movie idx={movie.id} title={movie.title}/>)

  
  return (
    <>
      <h1>Movie Recommender</h1>
      <div className="card">
        <input type='text' placeholder='Start searching :)'
        onChange={e => setSearchData(e.target.value)}/>
        
        {movies}
      </div>

    </>
  )
}

export default App
