import { useEffect, useState } from 'react'
import './App.css'
import { Book } from './modules/Book/Book'

interface BookInterface {
  id: Number
  title: string
  img_path: string
}

interface SearchMessage{
  searchbar_input: string
  search_type: SearchType
}

type SearchType = "exact" | "semantic"

function App() {
  const [searchData, setSearchData] = useState("")
  const [requestInputData, setRequestInputData] = useState("")
  const [bookArray, setBookArray] = useState([])
  const [searchType, setSearchType] = useState<SearchType>("exact")
  const [responseFailedError, setResponseFailedError] = useState(false)


  // block request while user is typing
  useEffect(() => {
    setResponseFailedError(false)


    const timeOutId = setTimeout(() => setRequestInputData(searchData), 500)
    return () => clearTimeout(timeOutId)
  }, [searchData])

  useEffect(() => {
    setResponseFailedError(false)

    
    if (requestInputData != ""){ 

      const searchMessage: SearchMessage = { 
        searchbar_input: requestInputData,
        search_type: searchType
      }

      const requestOptions = {
        method: 'POST',
        body: JSON.stringify(searchMessage)
      }

      fetch('http://127.0.0.1:8080/search', requestOptions)
          .then(response => response.json())
          .then(data => {
            setBookArray(data['similar_books'] ? data['similar_books'] : [])
          })
          .catch(() => setResponseFailedError(true))
      
      }
  }, [requestInputData, searchType])


  const books = bookArray.map((book:BookInterface) => <Book key={book.id.toString()} idx={book.id} title={book.title} img_path={book.img_path}/>)
  
  return (
    <>
      <h1>Book Recommender</h1>
      <div className="card">
        {responseFailedError && <p>Error occured</p>}
        <div>
        <input type='text' placeholder='Start searching :)'
        onChange={e => setSearchData(e.target.value)}/>
        <button className={searchType == "exact" ? "selected-search-type" : ""}
        onClick={() => setSearchType("exact")}>EXACT</button>
        <button className={searchType == "semantic" ? "selected-search-type" : ""}
        onClick={() => setSearchType("semantic")}>SEMANTIC</button>
        </div>
        
        {books}
      </div>

    </>
  )
}

export default App
