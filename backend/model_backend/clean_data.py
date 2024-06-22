import pandas as pd


IMG_PATH = "http://image.tmdb.org/t/p/w185"


movies = pd.read_csv('./movies_metadata.csv')
movies['release_date'] = pd.to_datetime(movies['release_date'], errors='coerce')


wiki_plots = pd.read_csv('./wiki_movie_plots_deduped.csv')

movies['year'] = movies['release_date'].dt.year
movie_dataset = pd.merge(wiki_plots, movies, how='inner', left_on=['Title', 'Release Year'], right_on=['title', 'year'])


print(movie_dataset.columns)

print(movie_dataset['Plot'].apply(lambda pt: len(pt) if isinstance(pt, str) else 0).max())

print(movie_dataset[['year', 'original_title', 'Plot', 'poster_path']])


movie_dataset.to_csv('./movie_dataset.csv', index=False)