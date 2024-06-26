import './Book.css'
interface BookProps {
    idx: Number
    title: string
    img_path: string
}

const COVER_IMG_PREFIX = 'https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/'

export const Book: React.FC<BookProps> = ({idx, title, img_path}) => {
    return (
        <div className='movie-component' onClick={() => console.log(idx)}>
            <h3>{title}</h3>
            <img src={COVER_IMG_PREFIX+img_path} alt='' loading='lazy' />
        </div>
    )
}