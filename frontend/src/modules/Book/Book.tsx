import { useImage } from 'react-image'
import './Book.css'
import { useEffect, useState } from 'react'

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
    const hasError = useErrorBoundary()
    const {src} = useImage({
        'srcList': COVER_IMG_PREFIX+img_path
    })


    return !hasError ? <img src={src} alt=''/> : <div>Couldn't load an image</div>
}

const useErrorBoundary = () => {
    const [hasError, setHasError] = useState(false)
    useEffect(() => {
        const errorHandler = (error:ErrorEvent) => {
        console.error('Error caught by boundary:', error);
        setHasError(true);
        };

        // Register error handler
        window.addEventListener('error', errorHandler);

        // Clean up
        return () => {
        window.removeEventListener('error', errorHandler);
        };
    }, []);
    return hasError
}


export const Book: React.FC<BookProps> = ({idx, title, img_path, findSimilarBooks}) => {

    return (
            <div className='movie-component' onClick={() => findSimilarBooks(idx, title)}>
                <h3>{title}</h3>
                <BookCover img_path={img_path}/>
            </div>
    )
}