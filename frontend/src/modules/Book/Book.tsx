import { useImage } from 'react-image'
import './Book.css'

import { ErrorBoundary } from 'react-error-boundary'

interface BookProps {
    idx: number
    title: string
    img_path: string
    findSimilarBooks: (id:number, title:string) => void
}

interface BookCoverProps {
    img_path: string
}



const COVER_IMG_PREFIX = 'https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/'

const BookCover: React.FC<BookCoverProps> = ({img_path}) =>{

    const {src} = useImage({
        'srcList': COVER_IMG_PREFIX+img_path
    })

    return <img className='book-cover' src={src} alt=''/> 
}


export const Book: React.FC<BookProps> = ({idx, title, img_path, findSimilarBooks}) => {

    return (
            <div className='movie-component' onClick={() => findSimilarBooks(idx, title)}>
                <ErrorBoundary fallback={<div>Couldn't load an image</div>}>
                    {img_path != "NaN" ?  <BookCover img_path={img_path}/> : <div>Couldn't load an image</div>}
                </ErrorBoundary>
            </div>
    )
}

